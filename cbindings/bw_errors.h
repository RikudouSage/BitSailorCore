#ifndef BITWARDEN_ERRORS
#define BITWARDEN_ERRORS

#include <stdlib.h>
#include <string.h>

static __thread char* bitwarden_last_error;

static void bitwarden_clear_last_error(void) {
	if (bitwarden_last_error) {
		free(bitwarden_last_error);
		bitwarden_last_error = NULL;
	}
}

static void bitwarden_set_last_error_copy(const char* msg) {
	bitwarden_clear_last_error();
	if (!msg) {
		return;
	}

	size_t len = strlen(msg);
	char* copy = (char*)malloc(len + 1);
	if (!copy) {
		return;
	}
	memcpy(copy, msg, len + 1);
	bitwarden_last_error = copy;
}

static size_t bitwarden_get_last_error(char* buf, size_t buf_len) {
	if (!bitwarden_last_error) {
		if (buf && buf_len > 0) {
			buf[0] = '\0';
		}
		return 1;
	}

	size_t len = strlen(bitwarden_last_error) + 1;
	if (buf && buf_len > 0) {
		size_t to_copy = len <= buf_len ? len : buf_len - 1;
		memcpy(buf, bitwarden_last_error, to_copy);
		buf[to_copy] = '\0';
	}
	return len;
}

#endif