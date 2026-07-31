#include <stdio.h>

#ifndef ROLE
#define ROLE "unspecified"
#endif

#ifndef BUILD
#define BUILD "unspecified"
#endif

int main(void) {
    printf("role=%s build=%s\n", ROLE, BUILD);
    return 0;
}
