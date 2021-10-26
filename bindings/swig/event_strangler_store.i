%newobject EventStranglerStoreExists;
%extend EventStranglerStoreExists_return {
  ~EventStranglerStoreExists_return() {
    free($self->r1);
  }
}

%newobject EventStranglerStoreGet;
%extend EventStranglerStoreGet_return {
  ~EventStranglerStoreGet_return() {
    free($self->r0);
    free($self->r1);
  }
}

%newobject EventStranglerStorePut;
%newobject EventStranglerStoreDelete;
%newobject EventStranglerStoreClose;
