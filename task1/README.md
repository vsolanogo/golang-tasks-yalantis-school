# Golang School

## Task 1

&quot;CSV Sorter&quot; is a CLI application that allows sorting of its input presented as CSV-text.

### Technical details

**Required features:**
1. The application runs as a CLI application.
2. It reads STDIN line by line. The end of the input is an empty line.
3. Each line is a sequence of alpha-numeric words, separated by commas - comma-separated values (CSV). The number of words (values) is the same in each line.
4. The application sorts all lines alphabetically by the first value in each line using the Tree Sort algorithm.
5. The application prints the result immediately in STDOUT, when the user ends to enter input text (presses &lt;Enter&gt; at a new line).
6. The application supports options:

| **Option usage** | **Meaning**                                                                                |
| ---              | ---                                                                                        |
| **-i file-name** | Use a file with the name **file-name** as an input.                                        |
| **-o file-name** | Use a file with the name **file-name** as an output.                                       |
| **-h**           | The first line is a header that must be ignored during sorting but included in the output. |
| **-f N**         | Sort input lines by value number **N** (starts from 1).                                                    |
| **-r**           | Sort input lines in reverse order.                                                         |

**Optional features** (not required but appreciated):

1. Add the ability to use a second algorithm for sorting based on the Red-black tree. Accordingly, add one more option **-a** with possible values **1** or **2**, which chooses Tree Sort or Red-black tree algorithm to use. By default, the application uses the Tree Sort algorithm.
