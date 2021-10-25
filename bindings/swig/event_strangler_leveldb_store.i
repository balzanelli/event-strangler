%newobject strangler_leveldb_store_new;

%extend strangler_leveldb_store_new_return {
  ~strangler_leveldb_store_new_return() {
    free($self->r1);
  }
}
