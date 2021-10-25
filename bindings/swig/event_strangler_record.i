%extend strangler_record {
  ~strangler_record() {
    free($self->hash_key);
    free($self->status);
    free($self->created_at);
    free($self->expires_at);
    free($self);
  }
};
