%newobject EventStranglerDynamoDBStoreNew;

%extend EventStranglerDynamoDBStoreNew_return {
    ~EventStranglerDynamoDBStoreNew_return() {
      free($self->r1);
    }
}
