#include <stdio.h>
#include <signal.h>
#include <time.h>

void sigint(int signal) {
  printf("got milk\n");
}

int main() {
  signal(SIGINT, &sigint);

  struct timespec ts;
  ts.tv_sec = 30;
  // nanosleep returns EINTR when interrupted by a signal. Don't restart it.
  nanosleep(&ts, NULL);
  printf("got interrupted\n");
  return 0;
}
