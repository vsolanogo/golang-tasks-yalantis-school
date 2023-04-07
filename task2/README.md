# Golang School

## Task 2

&quot;CSV Concurrent Sorter&quot; is a CLI application that allows sorting of its input presented as CSV-text.

### Technical details

Using the &quot;CSV Sorter&quot; from the [Task 1](https://github.com/oreuta/golang-school-task1), extend it with the following **required features:**

1. The application has additional option **-d dir-name** that specifies a directory where it must read input files from. All files in the directory must have the same format. The output stays the same, it is a one file with sorted content from all input files.
2. Processing must be implemented concurrently based on pipeline. The pipeline includes two stages:
  - Reading - read input and sent it line by line further.
  - Sorting - add received lines into the Tree.
3. The application outputs the result when the input ends up.
4. The project includes Unit tests covering the unit that builds the Tree.

**Optional features** (not required but appreciated):

1. Add signal processing that allows to gracefully stop the application when the user interrupts it pressing Ctrl-C. The interrupted application must write the current result.
2. If your application supports two types of algorithms, include to the project benchmarks comparing usages of these algorithms.
