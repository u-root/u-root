pkg/core is a set of packages and a struct that support building a (very)
limited subset of u-root commands into programs.

This was an experiment, of the kind golang occasionally runs.

Overall, it is not a good fit to u-root, as it results in packages that define
and parse flags and os.Args, among other things.  This can conflict with
programs using these packages.

Further, some pkg/core packages embed the name of the program in the package,
which made it confusing when error messages were printed. 
While this problem is easily fixed, it did indicate a flaw in the overall idea.

Programs don't belong in pkg/.

These packages are left here for now, but should not be used, extended, or added to,
beyond existing usage. They were an experiment that did not quite work out.

If you do want to embed u-root commands in your program, or bare metal system,
see pkg2cmd. It uses the same AST processing that u-root uses, and is hence
very much in line with our model. Further, it works in a side range of environemnts,
e.g. Plan 9 and bare metal.
