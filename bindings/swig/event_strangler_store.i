%newobject strangler_store_exists;
%extend strangler_store_exists_return {
  ~strangler_store_exists_return() {
    free($self->r1);
  }
}

%newobject strangler_store_get;
%extend strangler_store_get_return {
  ~strangler_store_get_return() {
    free($self->r0);
    free($self->r1);
  }
}

%newobject strangler_store_put;
%newobject strangler_store_delete;
%newobject strangler_store_close;
