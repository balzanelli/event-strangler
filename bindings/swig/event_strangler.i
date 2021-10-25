%{
#include "event_strangler.h"
%}
%include <stdint.i>

%ignore GoString;
%ignore _GoString_;
%ignore GoInterface;
%ignore GoSlice;

%include "event_strangler_lib.i"
%include "event_strangler_leveldb_store.i"
%include "event_strangler_record.i"
%include "event_strangler_store.i"

%newobject strangler_new;
%newobject strangler_complete;
%newobject strangler_purge;

%extend strangler_new_return {
  ~strangler_new_return() {
    free($self->r1);
  }
}
