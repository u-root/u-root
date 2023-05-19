#include <stdio.h>
#include <sys/types.h>
#include <unistd.h>

int main() {
  pid_t pid = fork();
  if (pid == 0) {
    printf("child, running execve\n");
    char *newargv[] = { "hello", NULL };
    char *newenviron[] = { NULL };
    execve("./hello", newargv, newenviron);
  } else {
    printf("parent\n");
  }
  return 0;
}
