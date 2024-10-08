#include "textflag.h"

TEXT ·BytesToNumericBytes(SB), NOSPLIT, $0
    MOVQ    len+8(FP), CX            // Load the length of the byte slice into CX
    MOVQ    b+0(FP), DX              // Load the pointer to the byte slice into DX
    MOVQ    DX, SI                   // SI will be our destination pointer for valid bytes
    XORQ    AX, AX                   // Clear AX to indicate no errors

    LEAQ    -1(DX)(CX*1), DI         // DI = DX + CX - 1 (point to the last byte of the array)

    MOVQ    $0, BX                   // Set initial state to q0

main_loop:
    CMPQ    CX, $0                   // Check if the length is 0
    JLE     accept                   // If it is, we are done
    

    // Load the next state based on the value in BX, ordered by the most common used states
    CMPQ    BX, $0                   // If BX == 0, jump to q0
    JE      q0
    CMPQ    BX, $1                   // If BX == 1, jump to q1
    JE      q1
    CMPQ    BX, $4                   // If BX == 4, jump to q3
    JE      q4
    CMPQ    BX, $5                   // If BX == 5, jump to q4
    JE      q5
    CMPQ    BX, $3                   // If BX == 3, jump to q3
    JE      q3
    CMPQ    BX, $2                   // If BX == 2, jump to q2
    JE      q2
    CMPQ    BX, $6                   // If BX == 6, jump to q6
    JE      q6

    JMP     error_invalid_state      // If we get here, we are in an invalid state

// State q0
q0:
    MOVB    (DX), AL                 // Load the current byte into AL

    CMPB    AL, $'1'                 // Check if AL is '1'
    JL      q0_less_than_one         // If less than '1', jump to q0_less_than_one

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JG      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q1              // Process the character

    q0_less_than_one:
    CMPB    AL, $'0'                 // Check if AL is '0'
    JNE     q0_char_sign                // If less than '0', jump to less_than_zero

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q3             // Jump to q3

    q0_char_sign:
    CMPB    AL, $'-'                 // Check if AL is '-'
    JE      set_state_q2             // If equal, jump to q2

    CMPB    AL, $'+'                 // Check if AL is '+'
    JE      set_state_q2             // If equal, jump to q2

    
    JMP     skip_char                // If we get here, we are in an invalid state

    set_state_q2:
    // save the sign character 
    MOVB    AL, (SI)                 // Store the sign character
    MOVQ    $2, BX                   // Set state to q2
    JMP     process                  // Process the character

q1:
    MOVB    (DX), AL                 // Load the current byte into AL

    CMPB    AL, $'0'                 // Check if AL is '0'
    JL      q1_less_than_zero        // If less than '0', jump to less_than_zero

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JG      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q1              // Process the character

    q1_less_than_zero:
    CMPB    AL, $'.'                 // Check if AL is '.'
    JE      set_state_q4             // If equal, jump to q4

    CMPB    AL, $' '                 // If we get here, we are in an invalid state
    JE      skip_char                // If equal, jump to error

set_state_q1:
    MOVQ    $1, BX                   // Set state to q1
    JMP     process                  // Process the character

set_state_q4:
    MOVQ    $4, BX                   // Set state to q4
    JMP     process                  // Process the '.' character

q2:
    MOVB    (DX), AL                 // Load the current byte into AL

    CMPB    AL, $'0'                 // Check if AL is '0'
    JL      error_unexpected_decimal // If less than '0', jump to error

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JG      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value

    // If the value is 0, jmp to q3, otherwise, jmp to q1 (you already changed the value of AL)
    CMPB    AL, $'0'                 // Check if AL is '0'
    JE      set_state_q3             // If equal, jump to q3

    JMP     set_state_q1             // Jump to q1

set_state_q3:
    MOVQ    $3, BX                   // Set state to q3
    JMP     process                  // Process the '.' character

q3:
    MOVB    (DX), AL                 // Load the current byte into AL

    CMPB    AL, $'.'                 // Check if AL is '.'  
    JNE     error_unexpected_decimal // If not, jump to error

    // save the decimal point
    MOVB    AL, (SI)                 // Store the decimal point

    MOVQ    $4, BX                   // Set state to q4
    JMP     process                  // Process the '.' character

q4:
    MOVB    (DX), AL                 // Load the current byte into AL

    CMPB    AL, $'1'                 // Check if AL is '1'
    JL      q4_less_than_one         // If less than '1', jump to q4_less_than_one

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JG      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q5             // Process the character

    q4_less_than_one:
    CMPB    AL, $'0'                 // Check if AL is '0'
    JNE     error_unexpected_decimal // If less than '0', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set state to q6
    JMP     process                  // Process the character

set_state_q5:
    MOVQ    $5, BX                   // Set state to q5
    JMP     process                  // Process the character

q5:
    MOVB    (DX), AL                 // Load the current byte into AL

    CMPB    AL, $'1'                 // Check if AL is '1'
    JL      q5_less_than_one         // If less than '1', jump to q5_less_than_one

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JG      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q5              // Process the character

    q5_less_than_one:
    CMPB    AL, $'0'                 // Check if AL is '0'
    JNE     skip_char                // If less than '0', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set state to q6
    JMP     process                  // Process the character

q6:
    MOVB    (DX), AL                 // Load the current byte into AL

    CMPB    AL, $'1'                 // Check if AL is '1'
    JL      q6_less_than_one         // If less than '1', jump to q6_less_than_one

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JG      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set state to q6
    JMP     process                  // Process the character

    q6_less_than_one:
    CMPB    AL, $'0'                 // Check if AL is '0' 
    JNE     error_unexpected_decimal // If less than '0', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q5             // Jump to q5

process:
    CMPQ    SI, DI                   // Check if SI (valid char pointer) reaches DI (invalid char pointer)
    JG     done                     // If SI >= DI, we're done
    MOVB    AL, (SI)                 // Store the numeric value
    INCQ    SI                       // Move the destination pointer
    INCQ    DX                       // Move the source pointer
    DECQ    CX                       // Decrement the length
    
    JMP     main_loop                // Continue the loop

skip_char:
    CMPQ    SI, DI                   // Check if SI reaches DI
    JG     done                     // If SI >= DI, we're done

    MOVB    $59, (DI)                // Store the invalid character
    DECQ    DI                       // Move the destination pointer
    INCQ    DX                       // Move the source pointer
    DECQ    CX                       // Decrement the length
    JMP     main_loop                // Continue the loop

accept:
    CMPQ    CX, $0                   // Check if the length is 0
    JNE     error_invalid_char       // If it is not, we are done

done:
    RET

error_invalid_char: // code -1 is 255 in uint8
    MOVQ    $-1, AX                  // Set error code for invalid character
    JMP     save_error

error_invalid_state: // code -2 is 254 in uint8
    MOVQ    $-2, AX                  // Set error code for invalid state
    JMP     save_error

error_unexpected_decimal: // code -3 is 253 in uint8
    MOVQ    $-3, AX                  // Set error code for unexpected decimal point
    JMP     save_error

error_multiple_decimals: // code -4 is 252 in uint8
    MOVQ    $-4, AX                  // Set error code for multiple decimal points
    JMP     save_error

save_error:
    MOVB    AL, (SI)                 // Store the invalid character
    INCQ    SI                       // Move the destination pointer
    DECQ    CX                       // Decrement the length
    JMP     accept
