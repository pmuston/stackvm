; Memory operations example
; Calculate c = a + b where a, b, c are in memory

# Memory layout:
# 0 = a (10)
# 1 = b (5)
# 2 = c (result)

START:
    PUSH 10         ; Initialize a
    STORE 0
    PUSH 5          ; Initialize b
    STORE 1

    LOAD 0          ; Load a
    LOAD 1          ; Load b
    ADD             ; a + b
    STORE 2         ; Store to c

    LOAD 2          ; Load result
    HALT
