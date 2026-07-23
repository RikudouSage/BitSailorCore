#ifndef BITWARDEN_GENERATOR
#define BITWARDEN_GENERATOR

#include <stdbool.h>

typedef struct {
    bool* lowercase;
    bool* uppercase;
    bool* numbers;
    bool* special;

    int* length;

    bool* avoidAmbiguous;
    int* minLowercase;
    int* minUppercase;
    int* minNumber;
    int* minSpecial;
} BitwardenPasswordGeneratorRequest;

typedef struct {
    int* numWords;
    char* wordSeparator;
    bool* capitalize;
    bool* includeNumber;
} BitwardenPassphraseGeneratorRequest;

#endif
