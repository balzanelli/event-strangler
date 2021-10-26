%extend EventStranglerRecord {
  ~EventStranglerRecord() {
    free($self->HashKey);
    free($self->Status);
    free($self->CreatedAt);
    free($self->ExpiresAt);
    free($self);
  }
};
