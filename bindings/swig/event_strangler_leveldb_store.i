%newobject EventStranglerLevelDBStoreNew;

%extend EventStranglerLevelDBStoreNew_return {
  ~EventStranglerLevelDBStoreNew_return() {
    free($self->r1);
  }
}
