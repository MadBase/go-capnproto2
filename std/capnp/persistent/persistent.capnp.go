// Code generated by capnpc-go. DO NOT EDIT.

package persistent

import (
	context "golang.org/x/net/context"
	capnp "github.com/MadBase/go-capnproto2/v2"
	text "github.com/MadBase/go-capnproto2/v2/encoding/text"
	schemas "github.com/MadBase/go-capnproto2/v2/schemas"
	server "github.com/MadBase/go-capnproto2/v2/server"
)

const PersistentAnnotation = uint64(0xf622595091cafb67)

type Persistent struct{ Client capnp.Client }

// Persistent_TypeID is the unique identifier for the type Persistent.
const Persistent_TypeID = 0xc8cb212fcd9f5691

func (c Persistent) Save(ctx context.Context, params func(Persistent_SaveParams) error, opts ...capnp.CallOption) Persistent_SaveResults_Promise {
	if c.Client == nil {
		return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0xc8cb212fcd9f5691,
			MethodID:      0,
			InterfaceName: "persistent.capnp:Persistent",
			MethodName:    "save",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Persistent_SaveParams{Struct: s}) }
	}
	return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Persistent_Server interface {
	Save(Persistent_save) error
}

func Persistent_ServerToClient(s Persistent_Server) Persistent {
	c, _ := s.(server.Closer)
	return Persistent{Client: server.New(Persistent_Methods(nil, s), c)}
}

func Persistent_Methods(methods []server.Method, s Persistent_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xc8cb212fcd9f5691,
			MethodID:      0,
			InterfaceName: "persistent.capnp:Persistent",
			MethodName:    "save",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Persistent_save{c, opts, Persistent_SaveParams{Struct: p}, Persistent_SaveResults{Struct: r}}
			return s.Save(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// Persistent_save holds the arguments for a server call to Persistent.save.
type Persistent_save struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Persistent_SaveParams
	Results Persistent_SaveResults
}

type Persistent_SaveParams struct{ capnp.Struct }

// Persistent_SaveParams_TypeID is the unique identifier for the type Persistent_SaveParams.
const Persistent_SaveParams_TypeID = 0xf76fba59183073a5

func NewPersistent_SaveParams(s *capnp.Segment) (Persistent_SaveParams, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveParams{st}, err
}

func NewRootPersistent_SaveParams(s *capnp.Segment) (Persistent_SaveParams, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveParams{st}, err
}

func ReadRootPersistent_SaveParams(msg *capnp.Message) (Persistent_SaveParams, error) {
	root, err := msg.RootPtr()
	return Persistent_SaveParams{root.Struct()}, err
}

func (s Persistent_SaveParams) String() string {
	str, _ := text.Marshal(0xf76fba59183073a5, s.Struct)
	return str
}

func (s Persistent_SaveParams) SealFor() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Persistent_SaveParams) HasSealFor() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Persistent_SaveParams) SealForPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Persistent_SaveParams) SetSealFor(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Persistent_SaveParams) SetSealForPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// Persistent_SaveParams_List is a list of Persistent_SaveParams.
type Persistent_SaveParams_List struct{ capnp.List }

// NewPersistent_SaveParams creates a new list of Persistent_SaveParams.
func NewPersistent_SaveParams_List(s *capnp.Segment, sz int32) (Persistent_SaveParams_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Persistent_SaveParams_List{l}, err
}

func (s Persistent_SaveParams_List) At(i int) Persistent_SaveParams {
	return Persistent_SaveParams{s.List.Struct(i)}
}

func (s Persistent_SaveParams_List) Set(i int, v Persistent_SaveParams) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s Persistent_SaveParams_List) String() string {
	str, _ := text.MarshalList(0xf76fba59183073a5, s.List)
	return str
}

// Persistent_SaveParams_Promise is a wrapper for a Persistent_SaveParams promised by a client call.
type Persistent_SaveParams_Promise struct{ *capnp.Pipeline }

func (p Persistent_SaveParams_Promise) Struct() (Persistent_SaveParams, error) {
	s, err := p.Pipeline.Struct()
	return Persistent_SaveParams{s}, err
}

func (p Persistent_SaveParams_Promise) SealFor() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type Persistent_SaveResults struct{ capnp.Struct }

// Persistent_SaveResults_TypeID is the unique identifier for the type Persistent_SaveResults.
const Persistent_SaveResults_TypeID = 0xb76848c18c40efbf

func NewPersistent_SaveResults(s *capnp.Segment) (Persistent_SaveResults, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveResults{st}, err
}

func NewRootPersistent_SaveResults(s *capnp.Segment) (Persistent_SaveResults, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveResults{st}, err
}

func ReadRootPersistent_SaveResults(msg *capnp.Message) (Persistent_SaveResults, error) {
	root, err := msg.RootPtr()
	return Persistent_SaveResults{root.Struct()}, err
}

func (s Persistent_SaveResults) String() string {
	str, _ := text.Marshal(0xb76848c18c40efbf, s.Struct)
	return str
}

func (s Persistent_SaveResults) SturdyRef() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Persistent_SaveResults) HasSturdyRef() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Persistent_SaveResults) SturdyRefPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Persistent_SaveResults) SetSturdyRef(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Persistent_SaveResults) SetSturdyRefPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// Persistent_SaveResults_List is a list of Persistent_SaveResults.
type Persistent_SaveResults_List struct{ capnp.List }

// NewPersistent_SaveResults creates a new list of Persistent_SaveResults.
func NewPersistent_SaveResults_List(s *capnp.Segment, sz int32) (Persistent_SaveResults_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Persistent_SaveResults_List{l}, err
}

func (s Persistent_SaveResults_List) At(i int) Persistent_SaveResults {
	return Persistent_SaveResults{s.List.Struct(i)}
}

func (s Persistent_SaveResults_List) Set(i int, v Persistent_SaveResults) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s Persistent_SaveResults_List) String() string {
	str, _ := text.MarshalList(0xb76848c18c40efbf, s.List)
	return str
}

// Persistent_SaveResults_Promise is a wrapper for a Persistent_SaveResults promised by a client call.
type Persistent_SaveResults_Promise struct{ *capnp.Pipeline }

func (p Persistent_SaveResults_Promise) Struct() (Persistent_SaveResults, error) {
	s, err := p.Pipeline.Struct()
	return Persistent_SaveResults{s}, err
}

func (p Persistent_SaveResults_Promise) SturdyRef() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type RealmGateway struct{ Client capnp.Client }

// RealmGateway_TypeID is the unique identifier for the type RealmGateway.
const RealmGateway_TypeID = 0x84ff286cd00a3ed4

func (c RealmGateway) Import(ctx context.Context, params func(RealmGateway_import_Params) error, opts ...capnp.CallOption) Persistent_SaveResults_Promise {
	if c.Client == nil {
		return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      0,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "import",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 2}
		call.ParamsFunc = func(s capnp.Struct) error { return params(RealmGateway_import_Params{Struct: s}) }
	}
	return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}
func (c RealmGateway) Export(ctx context.Context, params func(RealmGateway_export_Params) error, opts ...capnp.CallOption) Persistent_SaveResults_Promise {
	if c.Client == nil {
		return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      1,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "export",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 2}
		call.ParamsFunc = func(s capnp.Struct) error { return params(RealmGateway_export_Params{Struct: s}) }
	}
	return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type RealmGateway_Server interface {
	Import(RealmGateway_import) error

	Export(RealmGateway_export) error
}

func RealmGateway_ServerToClient(s RealmGateway_Server) RealmGateway {
	c, _ := s.(server.Closer)
	return RealmGateway{Client: server.New(RealmGateway_Methods(nil, s), c)}
}

func RealmGateway_Methods(methods []server.Method, s RealmGateway_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 2)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      0,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "import",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := RealmGateway_import{c, opts, RealmGateway_import_Params{Struct: p}, Persistent_SaveResults{Struct: r}}
			return s.Import(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      1,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "export",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := RealmGateway_export{c, opts, RealmGateway_export_Params{Struct: p}, Persistent_SaveResults{Struct: r}}
			return s.Export(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// RealmGateway_import holds the arguments for a server call to RealmGateway.import.
type RealmGateway_import struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  RealmGateway_import_Params
	Results Persistent_SaveResults
}

// RealmGateway_export holds the arguments for a server call to RealmGateway.export.
type RealmGateway_export struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  RealmGateway_export_Params
	Results Persistent_SaveResults
}

type RealmGateway_import_Params struct{ capnp.Struct }

// RealmGateway_import_Params_TypeID is the unique identifier for the type RealmGateway_import_Params.
const RealmGateway_import_Params_TypeID = 0xf0c2cc1d3909574d

func NewRealmGateway_import_Params(s *capnp.Segment) (RealmGateway_import_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_import_Params{st}, err
}

func NewRootRealmGateway_import_Params(s *capnp.Segment) (RealmGateway_import_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_import_Params{st}, err
}

func ReadRootRealmGateway_import_Params(msg *capnp.Message) (RealmGateway_import_Params, error) {
	root, err := msg.RootPtr()
	return RealmGateway_import_Params{root.Struct()}, err
}

func (s RealmGateway_import_Params) String() string {
	str, _ := text.Marshal(0xf0c2cc1d3909574d, s.Struct)
	return str
}

func (s RealmGateway_import_Params) Cap() Persistent {
	p, _ := s.Struct.Ptr(0)
	return Persistent{Client: p.Interface().Client()}
}

func (s RealmGateway_import_Params) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RealmGateway_import_Params) SetCap(v Persistent) error {
	if v.Client == nil {
		return s.Struct.SetPtr(0, capnp.Ptr{})
	}
	seg := s.Segment()
	in := capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	return s.Struct.SetPtr(0, in.ToPtr())
}

func (s RealmGateway_import_Params) Params() (Persistent_SaveParams, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	return Persistent_SaveParams{Struct: p.Struct()}, err
}

func (s RealmGateway_import_Params) HasParams() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s RealmGateway_import_Params) SetParams(v Persistent_SaveParams) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewParams sets the params field to a newly
// allocated Persistent_SaveParams struct, preferring placement in s's segment.
func (s RealmGateway_import_Params) NewParams() (Persistent_SaveParams, error) {
	ss, err := NewPersistent_SaveParams(s.Struct.Segment())
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// RealmGateway_import_Params_List is a list of RealmGateway_import_Params.
type RealmGateway_import_Params_List struct{ capnp.List }

// NewRealmGateway_import_Params creates a new list of RealmGateway_import_Params.
func NewRealmGateway_import_Params_List(s *capnp.Segment, sz int32) (RealmGateway_import_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return RealmGateway_import_Params_List{l}, err
}

func (s RealmGateway_import_Params_List) At(i int) RealmGateway_import_Params {
	return RealmGateway_import_Params{s.List.Struct(i)}
}

func (s RealmGateway_import_Params_List) Set(i int, v RealmGateway_import_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s RealmGateway_import_Params_List) String() string {
	str, _ := text.MarshalList(0xf0c2cc1d3909574d, s.List)
	return str
}

// RealmGateway_import_Params_Promise is a wrapper for a RealmGateway_import_Params promised by a client call.
type RealmGateway_import_Params_Promise struct{ *capnp.Pipeline }

func (p RealmGateway_import_Params_Promise) Struct() (RealmGateway_import_Params, error) {
	s, err := p.Pipeline.Struct()
	return RealmGateway_import_Params{s}, err
}

func (p RealmGateway_import_Params_Promise) Cap() Persistent {
	return Persistent{Client: p.Pipeline.GetPipeline(0).Client()}
}

func (p RealmGateway_import_Params_Promise) Params() Persistent_SaveParams_Promise {
	return Persistent_SaveParams_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

type RealmGateway_export_Params struct{ capnp.Struct }

// RealmGateway_export_Params_TypeID is the unique identifier for the type RealmGateway_export_Params.
const RealmGateway_export_Params_TypeID = 0xecafa18b482da3aa

func NewRealmGateway_export_Params(s *capnp.Segment) (RealmGateway_export_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_export_Params{st}, err
}

func NewRootRealmGateway_export_Params(s *capnp.Segment) (RealmGateway_export_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_export_Params{st}, err
}

func ReadRootRealmGateway_export_Params(msg *capnp.Message) (RealmGateway_export_Params, error) {
	root, err := msg.RootPtr()
	return RealmGateway_export_Params{root.Struct()}, err
}

func (s RealmGateway_export_Params) String() string {
	str, _ := text.Marshal(0xecafa18b482da3aa, s.Struct)
	return str
}

func (s RealmGateway_export_Params) Cap() Persistent {
	p, _ := s.Struct.Ptr(0)
	return Persistent{Client: p.Interface().Client()}
}

func (s RealmGateway_export_Params) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RealmGateway_export_Params) SetCap(v Persistent) error {
	if v.Client == nil {
		return s.Struct.SetPtr(0, capnp.Ptr{})
	}
	seg := s.Segment()
	in := capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	return s.Struct.SetPtr(0, in.ToPtr())
}

func (s RealmGateway_export_Params) Params() (Persistent_SaveParams, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	return Persistent_SaveParams{Struct: p.Struct()}, err
}

func (s RealmGateway_export_Params) HasParams() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s RealmGateway_export_Params) SetParams(v Persistent_SaveParams) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewParams sets the params field to a newly
// allocated Persistent_SaveParams struct, preferring placement in s's segment.
func (s RealmGateway_export_Params) NewParams() (Persistent_SaveParams, error) {
	ss, err := NewPersistent_SaveParams(s.Struct.Segment())
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// RealmGateway_export_Params_List is a list of RealmGateway_export_Params.
type RealmGateway_export_Params_List struct{ capnp.List }

// NewRealmGateway_export_Params creates a new list of RealmGateway_export_Params.
func NewRealmGateway_export_Params_List(s *capnp.Segment, sz int32) (RealmGateway_export_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return RealmGateway_export_Params_List{l}, err
}

func (s RealmGateway_export_Params_List) At(i int) RealmGateway_export_Params {
	return RealmGateway_export_Params{s.List.Struct(i)}
}

func (s RealmGateway_export_Params_List) Set(i int, v RealmGateway_export_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s RealmGateway_export_Params_List) String() string {
	str, _ := text.MarshalList(0xecafa18b482da3aa, s.List)
	return str
}

// RealmGateway_export_Params_Promise is a wrapper for a RealmGateway_export_Params promised by a client call.
type RealmGateway_export_Params_Promise struct{ *capnp.Pipeline }

func (p RealmGateway_export_Params_Promise) Struct() (RealmGateway_export_Params, error) {
	s, err := p.Pipeline.Struct()
	return RealmGateway_export_Params{s}, err
}

func (p RealmGateway_export_Params_Promise) Cap() Persistent {
	return Persistent{Client: p.Pipeline.GetPipeline(0).Client()}
}

func (p RealmGateway_export_Params_Promise) Params() Persistent_SaveParams_Promise {
	return Persistent_SaveParams_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

const schema_b8630836983feed7 = "x\xda\xbcSKh$E\x18\xfe\xbf\xean{F\x12" +
	"25\x9d\x90\x08\x89\xa3!\x12#N\x1e\x0a\x82A\x9c" +
	"\xe9\x80&9\x88\xd3\x13T\xe2\xc9\x9a\xd8\xd1\xc0LO" +
	"\xd3\xddy\x9dTP<\xa8\x87\xdc\xbcI\x08^r\xf1" +
	"\xe8\xe3\"\x06\x05\xf1\x01\x0a\x82gs\x13\x8c\xbb\xb3\xb0" +
	"\xbb\x87e\xa9\xa5z^\xbd\xc9\x84\xb0d\xd9KC\xff" +
	"\xf5W}\xff\xff=f\x0fPdsF\x9d\x119\xa3" +
	"\xc6C\xf2\xaf\x17\x1f\xfe\xa3\xfa\xa4\xfc\x908\xd7\xe4\xdf" +
	"\xff\x17>\x7f.\xb5\xf6\x0d\x11e`\x8d\xe1\x9a5\x05" +
	"\x93\xc8z\x02\x1f[yfZy6)\xbf\xbfR\xfc" +
	"\xf4\x87\xa5w\xbf&>\x0a\xb9\xf7\xfa\x17\xbf\xcf<\xfe" +
	"\xeb\xcfd\xc0\xcc\xe0\xd9e\xb6\x00k\x95\xa9+\xaf\xb1" +
	"\x02%\xce\xcf\xbe\xbe\xcb\x8e\xad\x8f\xd8$\x91\xb5\xcf\x16" +
	"\xad\x063\xad\x06\x1b\x96\x87\x07\xf9\xa5O\xf6\xbf\xfa\x8f" +
	"\xf8\xa3 2\x98z\xb5\xc1* X\xb7\xd96A\xbe" +
	"\xf2F\xfa\xf9\xb1\xdf\x8e\xae&\x1b\xfe\xd5\xe2\x86\xeb\x9a" +
	"jx\xe7\xd6/{\xa5\xd5\xf1\x1b\xf4'7\x1eC\x02" +
	"\x14\xd6?\xfa\xb1u\xa2\x9b\xd6\x89\x9e[\x1944\x10" +
	"\xe4\x97\xe1\xec\xc8\xeaw\xf5\x9b\xbd\xb61\x8cyXC" +
	"\x86\xda\x86\x1bj\x1b\xdf\x0d\xc2\x8d0r\x99\x17M\xaf" +
	"\x09\xdf\xf3\xe7\xcb\xae\xa8\xd6\x16E.r\xb7\xc5n\x09" +
	"pR\x9aA\xd4\x19\x12m\xb2\xf8\xdc<\x91\xfd4\xec" +
	"\x17\xc0?0\x81\xce\x9e\xdd\x8eM\xd5\xe1\xc3~\x1f\xfc" +
	"'\xb3\xb0Q\xf3\xebA\xc4\x91st\x86\xaeH\x00Q" +
	"\xbb\xd8\x99V\x95&\x9d\x14\x10\xe3\xaboF\x03\xb2H" +
	"\xdc#B\xd6\x00K\x16\x8a(\xb8;\x97\x061\xce\xa2" +
	"hw\xa3\x94\x00;\x03\xde_\xe1\xbc\xc2\x87\x02\xfeH" +
	" \x97\xbd\xc8\x0d<Q%\xb3\xec\xae\xcb\x97v\x92\x7f" +
	"\x9d\xb3\xdc\xab\xdb\x9e\x1btO[\xffm\x0d\xf4\x8e\x06" +
	"\xa5V\xc5\x8b\xa6W\xc4\x96[v\xc3\xcdj\x14\x92R" +
	"C\xd7t\"]\xad\xd3_&r\xfa48#\x0c2" +
	"\x8c6\x83\xb7w\xcb.a=\xa6)\xb1$\xb2=u" +
	"na\x98\xae\x17)\x0e\x12\xb6I\xbf\x99HD\xba\"" +
	"\xd5\x08%\x11\x08\xd2j\xa1l\xcfCf5\x0a\x1d=" +
	"\xf6F\xfbjWy\xfe\x14\x91\xdd\x07{\x14<o\x0e" +
	"\x84b\xcb=\xc3~OIT\xb1\xc9o\x0a\xdc(\xf3" +
	"\xf43r\xa5\xbb\xd9\xb9|\xb5<\x1b[v\xbai\x81" +
	"\x89\x92\x08LQ\x0b\x9dT\x87\xb0\xa9q\"gB\x83" +
	"\xb3\xc3\xc0\x81\xc1x\x02\xe5Q\xc7\xd7\xe0\xfc\xc8`\xae" +
	"\x09\x1f<\xc9]\x11\xf7\xc9\x9a\xe0\x84\x82/\x02Q\x0b" +
	"\x91\xe9r}/\x08\x17\xf9\x12\x99\x84\xce\xe7p\xd3\xcc" +
	"\xe0\x83\xe6\xe6\xc2\xc9/O\xceE\xf4gzf\xc0\xef" +
	"f \xce\x16\x98l|63\x9c}\xeb\xdb#R\xb0" +
	"\xf6\x08\xd0\xa7\x10\x0fe;\x92\xf0\"\xdb\xf3\xea\x91\x18" +
	"\x886\xea\x1e\xe9\x9dW\xb5\xf3\xd2[(\xc5\x9b\x9d\x0a" +
	"\xef\x02\x91\x9a\xdf\x19dx/tE\xf5\xe5z\xd0$" +
	"\xeaTr\xef\x04\x00\x00\xff\xff.7\xf3\xa1"

func init() {
	schemas.Register(schema_b8630836983feed7,
		0x84ff286cd00a3ed4,
		0xb76848c18c40efbf,
		0xc8cb212fcd9f5691,
		0xecafa18b482da3aa,
		0xf0c2cc1d3909574d,
		0xf622595091cafb67,
		0xf76fba59183073a5)
}
