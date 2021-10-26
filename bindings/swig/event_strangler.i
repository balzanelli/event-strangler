%{
#include "event_strangler.h"
%}
%include <stdint.i>

%ignore GoString;
%ignore _GoString_;
%ignore GoInterface;
%ignore GoSlice;

%include "event_strangler_lib.i"
%include "event_strangler_record.i"
%include "event_strangler_store.i"
%include "event_strangler_leveldb_store.i"

%newobject EventStranglerNew;
%newobject EventStranglerComplete;
%newobject EventStranglerPurge;

%extend EventStranglerNew_return {
  ~EventStranglerNew_return() {
    free($self->r1);
  }
}
