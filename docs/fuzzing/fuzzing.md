# Fuzzing

The goal of fuzzing is to continuously manipulate inputs to find bugs in our code. This comes with the advantage of having a vast input space not covered by the regular unit tests. The fuzzing tests are passed on a seed corpora, with each seed input being tested as if it were a unit test. Afterwards, the initial inputs are manipulated (transformed) and run again, over and over. Depending on the fuzzing engine, inputs are prioritized, for example, inputs that uncover new code paths (coverage-guided). It is only possible to find and cover a significant number of the edge-cases by writing fast fuzzing tests. These should be capable of 10-100,000 executions per second! The security of the tested code is asserted by running the fuzzing test for a long time without encountering any unexpected errors or panics. If any test fails, its corresponding inputs will be logged as a file so it can be reproduced.

## Golang

Fuzzing has been implemented in the Golang testing toolchain since 1.18. To get started on Golang native fuzzing, check out [Go Fuzz](https://tip.golang.org/security/fuzz/).
To run the fuzzing test FuzzFoo, run the command `go test -fuzz=FuzzFoo` inside of its directory. Additional arguments can be passed as well, e.g., `-fuzztime=30m` specifies the fuzzing tests to run for 30 minutes before passing. If no time is specified, the fuzzing will run indefinitely until aborted. If a test fails, the corresponding input is saved as a file under `testdata/fuzz/FuzzXXXX`. From now on, these inputs will be executed alongside the other seed inputs to make sure the test no longer fails for this particular input.

## Fuzzing guidelines for u-root:

- Write a fast fuzzing test. Execution time will heavily depend on your system, but at least 5K executions per second should be reachable even by most moderately fast systems.
- Pass enough seed inputs (by calling `f.Add(seed_input)`) to the fuzzing tests. Each input is executed every time by calling `go test` hence boosting coverage. The inputs also help the fuzzing engine understand what kind of format is being expected by the test, boosting its ability to find suitable inputs by manipulation.
- If a failing input is discovered, be sure to include it into the seed inputs of the test by hand and delete the input file. This is preferred as it reduces unnecessary clutter in the codebase.
- Think about what you want to fuzz, as fuzzing is not suitable in every test case. Focus on fuzzing functions which handle large amounts of untrusted data like parsers.
- Think about how you want to fuzz. Can you only detect crashes? Can you also test for any unexpected errors? Can you check if multiple parsing rounds of an object result in the same parsed object? This entirely depends on the functions under test and how the functions SHOULD behave.

## OSS-Fuzz

The [OSS-Fuzz project](https://google.github.io/oss-fuzz/) supports u-root by doing daily fuzzing runs for each of the fuzzing targets independently inside their own environment. Check out the [projects website](https://google.github.io/oss-fuzz/) for an overview of the architecture and supported tooling. At the moment, OSS-Fuzz does not support Golang's native fuzzing directly, so each fuzzing target is modified and run using [libfuzzer](https://www.llvm.org/docs/LibFuzzer.html) and [AddressSanitizer](https://github.com/google/sanitizers/wiki/AddressSanitizer) instead. Due to this, some modifications are necessary to make a seamless integration work. You can find the building steps for each of the fuzzing targets inside of the build.sh script inside of the u-root directory.

### Differences by fuzzing using libfuzzer

- Seed Corpora for libfuzzer are added by providing a single zip file containing all the individual seed files.
- Dictionaries can be used to enhance the fuzzing process. These contain tokens which are used frequently by the fuzzing engine by inserting them into the manipulated data. This is, for example, very beneficial in cases where a function is expecting unique values (magic values, addresses, etc.).
- The format of failed inputs of libfuzzer is not compatible with Golang's fuzzing (at least for now). This is NOT true for the input itself!

### Testing locally:

You can run the fuzzing targets in u-root via libfuzzer by using the [OSS-Fuzz repository](https://github.com/google/oss-fuzz).
Follow the steps below to build, run, and check the coverage of any fuzzing target. Be aware that the resulting coverage entirely depends on the number of unique inputs discovered.

### Running libfuzzer against a fuzzing target FuzzFoo in u-root

- Pull changes and build the dockerfile:
  `python3 infra/helper.py build_image u-root`
- Build all fuzzing targets:
  `python3 infra/helper.py build_fuzzers --sanitizer address u-root`
- Verify fuzzers are returning without errors:
  `python3 infra/helper.py check_build u-root`
- Run fuzzing target FuzzFoo i.e. fuzz_foo (check the build.sh for each ):
  `python3 infra/helper.py run_fuzzer u-root fuzz_foo`

### Collecting code coverage for a fuzzing target FuzzFo in u-root

- Pull coverage tooling:
  `python3 infra/helper.py pull_images`
- Run fuzzing target and collect corpora of unique inputs in a directory:

```
    mkdir -p /tmp/corpora
    python3 infra/helper.py run_fuzzer --corpus-dir=/tmp/corpora u-root fuzz_foo
```

- Build fuzzing target as coverage check:
  `python3 infra/helper.py build_fuzzers --sanitizer coverage u-root`
- Run for coverage with provided corpora:
  `python3 infra/helper.py coverage u-root --fuzz-target=fuzz_foo --corpus-dir=/tmp/corpora`
