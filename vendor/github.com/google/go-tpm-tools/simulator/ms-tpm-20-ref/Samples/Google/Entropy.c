#include <openssl/rand.h>

#include "Platform_fp.h"

// We get entropy from OpenSSL which gets its entropy from the OS.
int32_t _plat__GetEntropy(uint8_t *entropy, uint32_t amount) {
  if (RAND_bytes(entropy, amount) != 1) {
    return -1;
  }
  return amount;
}
