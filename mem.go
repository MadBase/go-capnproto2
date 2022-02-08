package capnp

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/MadBase/go-capnproto2/v2/internal/packed"
)

// Security limits. Matches C++ implementation.
const (
	defaultTraverseLimit = 64 << 20 // 64 MiB
	defaultDepthLimit    = 64

	maxStreamSegments = 512

	defaultDecodeLimit = 64 << 20 // 64 MiB
)

const maxDepth = ^uint(0)

// A Message is a tree of Cap'n Proto objects, split into one or more
// segments of contiguous memory.  The only required field is Arena.
// A Message is safe to read from multiple goroutines.
type Message struct {
	// rlimit must be first so that it is 64-bit aligned.
	// See sync/atomic docs.
	rlimit     ReadLimiter
	rlimitInit sync.Once

	Arena Arena

	// CapTable is the indexed list of the clients referenced in the
	// message.  Capability pointers inside the message will use this table
	// to map pointers to Clients.  The table is usually populated by the
	// RPC system.
	//
	// See https://capnproto.org/encoding.html#capabilities-interfaces for
	// more details on the capability table.
	CapTable []Client

	// TraverseLimit limits how many total bytes of data are allowed to be
	// traversed while reading.  Traversal is counted when a Struct or
	// List is obtained.  This means that calling a getter for the same
	// sub-struct multiple times will cause it to be double-counted.  Once
	// the traversal limit is reached, pointer accessors will report
	// errors. See https://capnproto.org/encoding.html#amplification-attack
	// for more details on this security measure.
	//
	// If not set, this defaults to 64 MiB.
	TraverseLimit uint64

	// DepthLimit limits how deeply-nested a message structure can be.
	// If not set, this defaults to 64.
	DepthLimit uint

	// mu protects the following fields:
	mu       sync.Mutex
	segs     map[SegmentID]*Segment
	firstSeg Segment // Preallocated first segment. msg is non-nil once initialized.
}

// NewMessage creates a message with a new root and returns the first
// segment.  It is an error to call NewMessage on an arena with data in it.
func NewMessage(arena Arena) (msg *Message, first *Segment, err error) {
	msg = &Message{Arena: arena}
	switch arena.NumSegments() {
	case 0:
		first, err = msg.allocSegment(wordSize)
		if err != nil {
			return nil, nil, err
		}
	case 1:
		first, err = msg.Segment(0)
		if err != nil {
			return nil, nil, err
		}
		if len(first.data) > 0 {
			return nil, nil, errHasData
		}
	default:
		return nil, nil, errHasData
	}
	if first.ID() != 0 {
		return nil, nil, errors.New("capnp: arena allocated first segment with non-zero ID")
	}
	seg, _, err := alloc(first, wordSize) // allocate root
	if err != nil {
		return nil, nil, err
	}
	if seg != first {
		return nil, nil, errors.New("capnp: arena didn't allocate first word in first segment")
	}
	return msg, first, nil
}

// Reset resets a message to use a different arena, allowing a single
// Message to be reused for reading multiple messages.  This invalidates
// any existing pointers in the Message, so use with caution.
func (m *Message) Reset(arena Arena) {
	m.mu.Lock()
	m.Arena = arena
	m.CapTable = nil
	m.segs = nil
	m.firstSeg = Segment{}
	m.mu.Unlock()
	if m.TraverseLimit == 0 {
		m.ReadLimiter().Reset(defaultTraverseLimit)
	} else {
		m.ReadLimiter().Reset(m.TraverseLimit)
	}
}

// Root returns the pointer to the message's root object.
//
// Deprecated: Use RootPtr.
func (m *Message) Root() (Pointer, error) {
	p, err := m.RootPtr()
	return p.toPointer(), err
}

// RootPtr returns the pointer to the message's root object.
func (m *Message) RootPtr() (Ptr, error) {
	s, err := m.Segment(0)
	if err != nil {
		return Ptr{}, err
	}
	return s.root().PtrAt(0)
}

// SetRoot sets the message's root object to p.
//
// Deprecated: Use SetRootPtr.
func (m *Message) SetRoot(p Pointer) error {
	return m.SetRootPtr(toPtr(p))
}

// SetRootPtr sets the message's root object to p.
func (m *Message) SetRootPtr(p Ptr) error {
	s, err := m.Segment(0)
	if err != nil {
		return err
	}
	return s.root().SetPtr(0, p)
}

// AddCap appends a capability to the message's capability table and
// returns its ID.
func (m *Message) AddCap(c Client) CapabilityID {
	n := CapabilityID(len(m.CapTable))
	m.CapTable = append(m.CapTable, c)
	return n
}

// ReadLimiter returns the message's read limiter.  Useful if you want
// to reset the traversal limit while reading.
func (m *Message) ReadLimiter() *ReadLimiter {
	m.rlimitInit.Do(func() {
		if m.TraverseLimit == 0 {
			m.rlimit.limit = defaultTraverseLimit
		} else {
			m.rlimit.limit = m.TraverseLimit
		}
	})
	return &m.rlimit
}

func (m *Message) depthLimit() uint {
	if m.DepthLimit != 0 {
		return m.DepthLimit
	}
	return defaultDepthLimit
}

// NumSegments returns the number of segments in the message.
func (m *Message) NumSegments() int64 {
	return int64(m.Arena.NumSegments())
}

// Segment returns the segment with the given ID.
func (m *Message) Segment(id SegmentID) (*Segment, error) {
	if isInt32Bit && id > maxInt32 {
		return nil, errSegment32Bit
	}
	if int64(id) >= m.Arena.NumSegments() {
		return nil, errSegmentOutOfBounds
	}
	m.mu.Lock()
	if seg := m.segment(id); seg != nil {
		m.mu.Unlock()
		return seg, nil
	}
	data, err := m.Arena.Data(id)
	if err != nil {
		m.mu.Unlock()
		return nil, err
	}
	seg := m.setSegment(id, data)
	m.mu.Unlock()
	return seg, nil
}

// segment returns the segment with the given ID.
// The caller must be holding m.mu.
func (m *Message) segment(id SegmentID) *Segment {
	if m.segs == nil {
		if id == 0 && m.firstSeg.msg != nil {
			return &m.firstSeg
		}
		return nil
	}
	return m.segs[id]
}

// setSegment creates or updates the Segment with the given ID.
// The caller must be holding m.mu.
func (m *Message) setSegment(id SegmentID, data []byte) *Segment {
	if m.segs == nil {
		if id == 0 {
			m.firstSeg = Segment{
				id:   id,
				msg:  m,
				data: data,
			}
			return &m.firstSeg
		}
		m.segs = make(map[SegmentID]*Segment)
		if m.firstSeg.msg != nil {
			m.segs[0] = &m.firstSeg
		}
	} else if seg := m.segs[id]; seg != nil {
		seg.data = data
		return seg
	}
	seg := &Segment{
		id:   id,
		msg:  m,
		data: data,
	}
	m.segs[id] = seg
	return seg
}

// allocSegment creates or resizes an existing segment such that
// cap(seg.Data) - len(seg.Data) >= sz.
func (m *Message) allocSegment(sz Size) (*Segment, error) {
	m.mu.Lock()
	if m.segs == nil && m.firstSeg.msg != nil {
		m.segs = make(map[SegmentID]*Segment)
		m.segs[0] = &m.firstSeg
	}
	id, data, err := m.Arena.Allocate(sz, m.segs)
	if err != nil {
		m.mu.Unlock()
		return nil, err
	}
	if isInt32Bit && id > maxInt32 {
		m.mu.Unlock()
		return nil, errSegment32Bit
	}
	seg := m.setSegment(id, data)
	m.mu.Unlock()
	return seg, nil
}

// alloc allocates sz zero-filled bytes.  It prefers using s, but may
// use a different segment in the same message if there's not sufficient
// capacity.
func alloc(s *Segment, sz Size) (*Segment, Address, error) {
	sz = sz.padToWord()
	if sz > maxSize-wordSize {
		return nil, 0, errOverflow
	}

	if !hasCapacity(s.data, sz) {
		var err error
		s, err = s.msg.allocSegment(sz)
		if err != nil {
			return nil, 0, err
		}
	}

	addr := Address(len(s.data))
	end, ok := addr.addSize(sz)
	if !ok {
		return nil, 0, errOverflow
	}
	space := s.data[len(s.data):end]
	s.data = s.data[:end]
	for i := range space {
		space[i] = 0
	}
	return s, addr, nil
}

// An Arena loads and allocates segments for a Message.
type Arena interface {
	// NumSegments returns the number of segments in the arena.
	// This must not be larger than 1<<32.
	NumSegments() int64

	// Data loads the data for the segment with the given ID.  IDs are in
	// the range [0, NumSegments()).
	// must be tightly packed in the range [0, NumSegments()).
	Data(id SegmentID) ([]byte, error)

	// Allocate selects a segment to place a new object in, creating a
	// segment or growing the capacity of a previously loaded segment if
	// necessary.  If Allocate does not return an error, then the
	// difference of the capacity and the length of the returned slice
	// must be at least minsz.  segs is a map of segment slices returned
	// by the Data method keyed by ID (although the length of these slices
	// may have changed by previous allocations).  Allocate must not
	// modify segs.
	//
	// If Allocate creates a new segment, the ID must be one larger than
	// the last segment's ID or zero if it is the first segment.
	//
	// If Allocate returns an previously loaded segment's ID, then the
	// arena is responsible for preserving the existing data in the
	// returned byte slice.
	Allocate(minsz Size, segs map[SegmentID]*Segment) (SegmentID, []byte, error)
}

type singleSegmentArena []byte

// SingleSegment returns a new arena with an expanding single-segment
// buffer.  b can be used to populate the segment for reading or to
// reserve memory of a specific size.  A SingleSegment arena does not
// return errors unless you attempt to access another segment.
func SingleSegment(b []byte) Arena {
	ssa := new(singleSegmentArena)
	*ssa = b
	return ssa
}

func (ssa *singleSegmentArena) NumSegments() int64 {
	return 1
}

func (ssa *singleSegmentArena) Data(id SegmentID) ([]byte, error) {
	if id != 0 {
		return nil, errSegmentOutOfBounds
	}
	return *ssa, nil
}

func (ssa *singleSegmentArena) Allocate(sz Size, segs map[SegmentID]*Segment) (SegmentID, []byte, error) {
	data := []byte(*ssa)
	if segs[0] != nil {
		data = segs[0].data
	}
	if len(data)%int(wordSize) != 0 {
		return 0, nil, errors.New("capnp: segment size is not a multiple of word size")
	}
	if hasCapacity(data, sz) {
		return 0, data, nil
	}
	inc, err := nextAlloc(int64(cap(data)), int64(maxSegmentSize()), sz)
	if err != nil {
		return 0, nil, fmt.Errorf("capnp: alloc %d bytes: %v", sz, err)
	}
	buf := make([]byte, len(data), cap(data)+inc)
	copy(buf, data)
	*ssa = buf
	return 0, *ssa, nil
}

type roSingleSegment []byte

func (ss roSingleSegment) NumSegments() int64 {
	return 1
}

func (ss roSingleSegment) Data(id SegmentID) ([]byte, error) {
	if id != 0 {
		return nil, errSegmentOutOfBounds
	}
	return ss, nil
}

func (ss roSingleSegment) Allocate(sz Size, segs map[SegmentID]*Segment) (SegmentID, []byte, error) {
	return 0, nil, errors.New("capnp: segment is read-only")
}

type multiSegmentArena [][]byte

// MultiSegment returns a new arena that allocates new segments when
// they are full.  b can be used to populate the buffer for reading or
// to reserve memory of a specific size.
func MultiSegment(b [][]byte) Arena {
	msa := new(multiSegmentArena)
	*msa = b
	return msa
}

// demuxArena slices b into a multi-segment arena.
func demuxArena(hdr streamHeader, data []byte) (Arena, error) {
	segs := make([][]byte, int(hdr.maxSegment())+1)
	for i := range segs {
		sz, err := hdr.segmentSize(uint32(i))
		if err != nil {
			return nil, err
		}
		segs[i], data = data[:sz:sz], data[sz:]
	}
	return MultiSegment(segs), nil
}

func (msa *multiSegmentArena) NumSegments() int64 {
	return int64(len(*msa))
}

func (msa *multiSegmentArena) Data(id SegmentID) ([]byte, error) {
	if int64(id) >= int64(len(*msa)) {
		return nil, errSegmentOutOfBounds
	}
	return (*msa)[id], nil
}

func (msa *multiSegmentArena) Allocate(sz Size, segs map[SegmentID]*Segment) (SegmentID, []byte, error) {
	var total int64
	for i, data := range *msa {
		id := SegmentID(i)
		if s := segs[id]; s != nil {
			data = s.data
		}
		if hasCapacity(data, sz) {
			return id, data, nil
		}
		total += int64(cap(data))
		if total < 0 {
			// Overflow.
			return 0, nil, fmt.Errorf("capnp: alloc %d bytes: message too large", sz)
		}
	}
	n, err := nextAlloc(total, 1<<63-1, sz)
	if err != nil {
		return 0, nil, fmt.Errorf("capnp: alloc %d bytes: %v", sz, err)
	}
	buf := make([]byte, 0, n)
	id := SegmentID(len(*msa))
	*msa = append(*msa, buf)
	return id, buf, nil
}

// nextAlloc computes how much more space to allocate given the number
// of bytes allocated in the entire message and the requested number of
// bytes.  It will always return a multiple of wordSize.  max must be a
// multiple of wordSize.  The sum of curr and the returned size will
// always be less than max.
func nextAlloc(curr, max int64, req Size) (int, error) {
	if req == 0 {
		return 0, nil
	}
	maxinc := int64(1<<32 - 8) // largest word-aligned Size
	if isInt32Bit {
		maxinc = 1<<31 - 8 // largest word-aligned int
	}
	if int64(req) > maxinc {
		return 0, errors.New("allocation too large")
	}
	req = req.padToWord()
	want := curr + int64(req)
	if want <= curr || want > max {
		return 0, errors.New("allocation overflows message size")
	}
	new := curr
	double := new + new
	switch {
	case want < 1024:
		next := (1024 - curr + 7) &^ 7
		if next < curr {
			return int((curr + 7) &^ 7), nil
		}
		return int(next), nil
	case want > double:
		return int(req), nil
	default:
		for 0 < new && new < want {
			new += new / 4
		}
		if new <= 0 {
			return int(req), nil
		}
		delta := new - curr
		if delta > maxinc {
			return int(maxinc), nil
		}
		return int((delta + 7) &^ 7), nil
	}
}

// A Decoder represents a framer that deserializes a particular Cap'n
// Proto input stream.
type Decoder struct {
	r io.Reader

	segbuf [msgHeaderSize]byte
	hdrbuf []byte

	reuse bool
	buf   []byte
	msg   Message
	arena roSingleSegment

	// Maximum number of bytes that can be read per call to Decode.
	// If not set, a reasonable default is used.
	MaxMessageSize uint64
}

// NewDecoder creates a new Cap'n Proto framer that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// NewPackedDecoder creates a new Cap'n Proto framer that reads from a
// packed stream r.
func NewPackedDecoder(r io.Reader) *Decoder {
	return NewDecoder(packed.NewReader(bufio.NewReader(r)))
}

// Decode reads a message from the decoder stream.
func (d *Decoder) Decode() (*Message, error) {
	maxSize := d.MaxMessageSize
	if maxSize == 0 {
		maxSize = defaultDecodeLimit
	}
	if _, err := io.ReadFull(d.r, d.segbuf[:]); err != nil {
		return nil, err
	}
	maxSeg := binary.LittleEndian.Uint32(d.segbuf[:])
	if maxSeg > maxStreamSegments {
		return nil, errTooManySegments
	}
	hdrSize := streamHeaderSize(maxSeg)
	if hdrSize > maxSize || hdrSize > (1<<31-1) {
		return nil, errDecodeLimit
	}
	d.hdrbuf = resizeSlice(d.hdrbuf, int(hdrSize))
	copy(d.hdrbuf, d.segbuf[:])
	if _, err := io.ReadFull(d.r, d.hdrbuf[msgHeaderSize:]); err != nil {
		return nil, err
	}
	hdr, _, err := parseStreamHeader(d.hdrbuf)
	if err != nil {
		return nil, err
	}
	total, err := hdr.totalSize()
	if err != nil {
		return nil, err
	}
	// TODO(someday): if total size is greater than can fit in one buffer,
	// attempt to allocate buffer per segment.
	if total > maxSize-hdrSize || total > (1<<31-1) {
		return nil, errDecodeLimit
	}
	if !d.reuse {
		buf := make([]byte, int(total))
		if _, err := io.ReadFull(d.r, buf); err != nil {
			return nil, err
		}
		arena, err := demuxArena(hdr, buf)
		if err != nil {
			return nil, err
		}
		return &Message{Arena: arena}, nil
	}
	d.buf = resizeSlice(d.buf, int(total))
	if _, err := io.ReadFull(d.r, d.buf); err != nil {
		return nil, err
	}
	var arena Arena
	if hdr.maxSegment() == 0 {
		d.arena = d.buf[:len(d.buf):len(d.buf)]
		arena = &d.arena
	} else {
		var err error
		arena, err = demuxArena(hdr, d.buf)
		if err != nil {
			return nil, err
		}
	}
	d.msg.Reset(arena)
	return &d.msg, nil
}

func resizeSlice(b []byte, size int) []byte {
	if cap(b) < size {
		return make([]byte, size)
	}
	return b[:size]
}

// ReuseBuffer causes the decoder to reuse its buffer on subsequent decodes.
// The decoder may return messages that cannot handle allocations.
func (d *Decoder) ReuseBuffer() {
	d.reuse = true
}

// Unmarshal reads an unpacked serialized stream into a message.  No
// copying is performed, so the objects in the returned message read
// directly from data.
func Unmarshal(data []byte) (*Message, error) {
	if len(data) == 0 {
		return nil, io.EOF
	}
	hdr, data, err := parseStreamHeader(data)
	if err != nil {
		return nil, err
	}
	if tot, err := hdr.totalSize(); err != nil {
		return nil, err
	} else if tot > uint64(len(data)) {
		return nil, io.ErrUnexpectedEOF
	}
	arena, err := demuxArena(hdr, data)
	if err != nil {
		return nil, err
	}
	return &Message{Arena: arena}, nil
}

// UnmarshalPacked reads a packed serialized stream into a message.
func UnmarshalPacked(data []byte) (*Message, error) {
	if len(data) == 0 {
		return nil, io.EOF
	}
	data, err := packed.Unpack(nil, data)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}

// MustUnmarshalRoot reads an unpacked serialized stream and returns
// its root pointer.  If there is any error, it panics.
//
// Deprecated: Use MustUnmarshalRootPtr.
func MustUnmarshalRoot(data []byte) Pointer {
	msg, err := Unmarshal(data)
	if err != nil {
		panic(err)
	}
	p, err := msg.Root()
	if err != nil {
		panic(err)
	}
	return p
}

// MustUnmarshalRootPtr reads an unpacked serialized stream and returns
// its root pointer.  If there is any error, it panics.
func MustUnmarshalRootPtr(data []byte) Ptr {
	msg, err := Unmarshal(data)
	if err != nil {
		panic(err)
	}
	p, err := msg.RootPtr()
	if err != nil {
		panic(err)
	}
	return p
}

// An Encoder represents a framer for serializing a particular Cap'n
// Proto stream.
type Encoder struct {
	w      io.Writer
	hdrbuf []byte
	bufs   [][]byte

	packed  bool
	packbuf []byte
}

// NewEncoder creates a new Cap'n Proto framer that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// NewPackedEncoder creates a new Cap'n Proto framer that writes to a
// packed stream w.
func NewPackedEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, packed: true}
}

// Encode writes a message to the encoder stream.
func (e *Encoder) Encode(m *Message) error {
	nsegs := m.NumSegments()
	if nsegs == 0 {
		return errMessageEmpty
	}
	e.bufs = append(e.bufs[:0], nil) // first element is placeholder for header
	maxSeg := uint32(nsegs - 1)
	hdrSize := streamHeaderSize(maxSeg)
	if uint64(cap(e.hdrbuf)) < hdrSize {
		e.hdrbuf = make([]byte, 0, hdrSize)
	}
	e.hdrbuf = appendUint32(e.hdrbuf[:0], maxSeg)
	for i := int64(0); i < nsegs; i++ {
		s, err := m.Segment(SegmentID(i))
		if err != nil {
			return err
		}
		n := len(s.data)
		if int64(n) > int64(maxSize) {
			return errSegmentTooLarge
		}
		e.hdrbuf = appendUint32(e.hdrbuf, uint32(Size(n)/wordSize))
		e.bufs = append(e.bufs, s.data)
	}
	if len(e.hdrbuf)%int(wordSize) != 0 {
		e.hdrbuf = appendUint32(e.hdrbuf, 0)
	}
	e.bufs[0] = e.hdrbuf
	if e.packed {
		return e.writePacked(e.bufs)
	}
	return e.write(e.bufs)
}

func (e *Encoder) writePacked(bufs [][]byte) error {
	for _, b := range bufs {
		e.packbuf = packed.Pack(e.packbuf[:0], b)
		if _, err := e.w.Write(e.packbuf); err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) segmentSizes() ([]Size, error) {
	nsegs := m.NumSegments()
	sizes := make([]Size, nsegs)
	for i := int64(0); i < nsegs; i++ {
		s, err := m.Segment(SegmentID(i))
		if err != nil {
			return sizes[:i], err
		}
		n := len(s.data)
		if int64(n) > int64(maxSize) {
			return sizes[:i], errSegmentTooLarge
		}
		sizes[i] = Size(n)
	}
	return sizes, nil
}

// Marshal concatenates the segments in the message into a single byte
// slice including framing.
func (m *Message) Marshal() ([]byte, error) {
	// Compute buffer size.
	// TODO(light): error out if too many segments
	nsegs := m.NumSegments()
	if nsegs == 0 {
		return nil, errMessageEmpty
	}
	maxSeg := uint32(nsegs - 1)
	hdrSize := streamHeaderSize(maxSeg)
	sizes, err := m.segmentSizes()
	if err != nil {
		return nil, err
	}
	// TODO(light): error out if too large
	total := uint64(hdrSize) + totalSize(sizes)

	// Fill in buffer.
	buf := make([]byte, hdrSize, total)
	// TODO: remove marshalStreamHeader and inline.
	marshalStreamHeader(buf, sizes)
	for i := int64(0); i < nsegs; i++ {
		s, err := m.Segment(SegmentID(i))
		if err != nil {
			return nil, err
		}
		buf = append(buf, s.data...)
	}
	return buf, nil
}

// MarshalPacked marshals the message in packed form.
func (m *Message) MarshalPacked() ([]byte, error) {
	data, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 0, len(data))
	buf = packed.Pack(buf, data)
	return buf, nil
}

// Stream header sizes.
const (
	msgHeaderSize = 4
	segHeaderSize = 4
)

// streamHeaderSize returns the size of the header, given the
// first 32-bit number.
func streamHeaderSize(n uint32) uint64 {
	return (msgHeaderSize + segHeaderSize*(uint64(n)+1) + 7) &^ 7
}

// marshalStreamHeader marshals the sizes into the byte slice, which
// must be of size streamHeaderSize(len(sizes) - 1).
//
// TODO: remove marshalStreamHeader and inline.
func marshalStreamHeader(b []byte, sizes []Size) {
	binary.LittleEndian.PutUint32(b, uint32(len(sizes)-1))
	for i, sz := range sizes {
		loc := msgHeaderSize + i*segHeaderSize
		binary.LittleEndian.PutUint32(b[loc:], uint32(sz/Size(wordSize)))
	}
}

// appendUint32 appends a uint32 to a byte slice and returns the
// new slice.
func appendUint32(b []byte, v uint32) []byte {
	b = append(b, 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(b[len(b)-4:], v)
	return b
}

type streamHeader struct {
	b []byte
}

// parseStreamHeader parses the header of the stream framing format.
func parseStreamHeader(data []byte) (h streamHeader, tail []byte, err error) {
	if uint64(len(data)) < streamHeaderSize(0) {
		return streamHeader{}, nil, io.ErrUnexpectedEOF
	}
	maxSeg := binary.LittleEndian.Uint32(data)
	// TODO(light): check int
	hdrSize := streamHeaderSize(maxSeg)
	if uint64(len(data)) < hdrSize {
		return streamHeader{}, nil, io.ErrUnexpectedEOF
	}
	return streamHeader{b: data}, data[hdrSize:], nil
}

func (h streamHeader) maxSegment() uint32 {
	return binary.LittleEndian.Uint32(h.b)
}

func (h streamHeader) segmentSize(i uint32) (Size, error) {
	s := binary.LittleEndian.Uint32(h.b[msgHeaderSize+i*segHeaderSize:])
	sz, ok := wordSize.times(int32(s))
	if !ok {
		return 0, errSegmentTooLarge
	}
	return sz, nil
}

func (h streamHeader) totalSize() (uint64, error) {
	var sum uint64
	for i := uint64(0); i <= uint64(h.maxSegment()); i++ {
		x, err := h.segmentSize(uint32(i))
		if err != nil {
			return sum, err
		}
		sum += uint64(x)
	}
	return sum, nil
}

func hasCapacity(b []byte, sz Size) bool {
	return sz <= Size(cap(b)-len(b))
}

func totalSize(s []Size) uint64 {
	var sum uint64
	for _, sz := range s {
		sum += uint64(sz)
	}
	return sum
}

const (
	maxInt32 = 0x7fffffff
	maxInt   = int(^uint(0) >> 1)

	isInt32Bit = maxInt == maxInt32
)

// maxSegmentSize returns the maximum permitted size of a single segment
// on this platform.
//
// This is effectively a compile-time constant, but can't be represented
// as a constant because it requires a conditional.  It is trivially
// inlinable and optimizable, so should act like one.
func maxSegmentSize() Size {
	if isInt32Bit {
		return Size(maxInt32 - 7)
	} else {
		return maxSize - 7
	}
}

var (
	errSegmentOutOfBounds = errors.New("capnp: segment ID out of bounds")
	errSegment32Bit       = errors.New("capnp: segment ID larger than 31 bits")
	errMessageEmpty       = errors.New("capnp: marshalling an empty message")
	errHasData            = errors.New("capnp: NewMessage called on arena with data")
	errSegmentTooLarge    = errors.New("capnp: segment too large")
	errTooManySegments    = errors.New("capnp: too many segments to decode")
	errDecodeLimit        = errors.New("capnp: message too large")
)
