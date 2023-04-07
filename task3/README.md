# Golang School

## Task 3

&quot;CSV Sorter System&quot; is a distributed application that allows sorting of its input presented as CSV-text.

### Technical details

Using the &quot;CSV Concurrent Sorter&quot; from the [Task 2](https://github.com/oreuta/golang-school-task2), build a system with the following **required features**:

1. The system includes two microservices MS1 and MS2 that communicate via gRPC.
2. MS1 reads and writes user data and implements UI part of &quot;CSV Concurrent Sorter&quot;.
3. MS2 sorts user data and implements sorting unit of &quot;CSV Concurrent Sorter&quot;.
4. Communication is based on requests from MS1 to MS2 and synchronous responses from MS2 to MS1.

Request is a massage with two fields **Action** and **Payload** :

| **Action** | **Payload**     | **Description**                                                     |
| ---        | ---             | ---                                                                 |
| START      | Header or Empty | MS2 is required to start building of a new Tree                     |
| ADD        | Array of Lines  | MS2 is required to add Lines to the existed Tree                    |
| GET        | Empty           | MS2 is required to return the current (sorted) content of the Tree  |
| STOP       | Empty           | MS2 is required to discard the Tree and be ready to build a new one |

Response is a message with two fields **Error** and **Payload**. **Error** contains error message or is empty. **Payload** is empty if error occurs or received **Action** was not GET. If the received action was GET, **Payload** in response message contains sorted output with header or not, depending on the content of the last request message with the **Action** START.

**Optional features** (not required but appreciated):

1. Add a request message with **Action** SHUTDOWN that requires MS2 to stop running and return response message as for Action GET request.
2. Add to message one more field **Algorithm** with values **1** or **2** that is not empty only in **Action** START request. It is used to choose sorting algorithm as it is described in Task 2.
