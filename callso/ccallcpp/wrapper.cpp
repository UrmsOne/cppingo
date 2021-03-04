#include "person.h"
#include "wrapper.h"

#ifdef __cplusplus
extern "C"{
#endif


void *call_Person_Create() {
  return new Person("urmsone", 18);
}

void call_Person_Destroy(void *p) {
  delete static_cast<Person *>(p);
}

int call_Person_GetAge(void *p) {
  return static_cast<Person *>(p)->GetAge();
}

const char *call_Person_GetName(void *p) {
  return static_cast<Person *>(p)->GetName();
}

#ifdef __cplusplus
}
#endif