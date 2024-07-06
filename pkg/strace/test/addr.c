#include <stdio.h>
#include <netdb.h>
#include <string.h>

int main() {
    struct addrinfo hints, *res;
    int status;

    memset(&hints, 0, sizeof hints);
    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;

    status = getaddrinfo("localhost", NULL, &hints, &res);
    if (status != 0) {
        return 1;
    }

    freeaddrinfo(res);
    return 0;
}
