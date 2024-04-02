from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    File: _ClassVar[QueryType]
    Link: _ClassVar[QueryType]

class ResponseStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    Success: _ClassVar[ResponseStatus]
    Processing: _ClassVar[ResponseStatus]
    Error: _ClassVar[ResponseStatus]
    Canceled: _ClassVar[ResponseStatus]
File: QueryType
Link: QueryType
Success: ResponseStatus
Processing: ResponseStatus
Error: ResponseStatus
Canceled: ResponseStatus

class Query(_message.Message):
    __slots__ = ("id", "type", "source", "timeout")
    ID_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_FIELD_NUMBER: _ClassVar[int]
    id: int
    type: QueryType
    source: str
    timeout: int
    def __init__(self, id: _Optional[int] = ..., type: _Optional[_Union[QueryType, str]] = ..., source: _Optional[str] = ..., timeout: _Optional[int] = ...) -> None: ...

class Response(_message.Message):
    __slots__ = ("status",)
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: ResponseStatus
    def __init__(self, status: _Optional[_Union[ResponseStatus, str]] = ...) -> None: ...

class ProcessedReq(_message.Message):
    __slots__ = ("queryId",)
    QUERYID_FIELD_NUMBER: _ClassVar[int]
    queryId: int
    def __init__(self, queryId: _Optional[int] = ...) -> None: ...

class ProcessedResp(_message.Message):
    __slots__ = ("src",)
    SRC_FIELD_NUMBER: _ClassVar[int]
    src: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, src: _Optional[_Iterable[str]] = ...) -> None: ...
