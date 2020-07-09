LIBNAME        = solwrap
SONAME         = lib$(LIBNAME).so

SOLWRAP_DIR    = $(CURDIR)/includes/solwrap
SOLCLIENT_DIR  = $(CURDIR)/includes/solclient

LIB_DIR        = $(CURDIR)
BUILD_DIR      = $(CURDIR)/bin

PYINC=
PYLIB=
PYDEF=

INCDIRS              = $(SOLCLIENT_DIR)/include $(SOLWRAP_DIR)/include
LIBDIRS              = $(CURDIR) $(SOLCLIENT_DIR)/lib
STATIC_LIBS          = $(CURDIR)/lib/libsolwrap.a $(SOLCLIENT_DIR)/lib/libsolclient.a
WRAP_LIBS            = pthread
WRAP_BIN_LIBS        = pthread rt
WRAP_GO_LIBS         = pthread rt stdc++
SYMS                 = PROVIDE_LOG_UTILITIES SOLCLIENT_CONST_PROPERTIES _REENTRANT _LINUX_X86_64
DEBUG                = -g

SOL_SRC     = $(wildcard $(SOLWRAP_DIR)/src/*.cpp)
SOL_TEST    = $(wildcard ./bin/*.test)

CXXFLAGS_LIB = $(foreach d, $(INCDIRS), -I$d) $(foreach s, $(SYMS), -D$s) -m64 $(DEBUG)
CXXFLAGS     = $(foreach d, $(INCDIRS), -I$d) $(foreach s, $(SYMS), -D$s) -m64 $(DEBUG)

WRAPPER_LIBS      = $(foreach l, $(WRAP_LIBS), -l$l)
WRAPPER_GO_LIBS   = $(foreach l, $(WRAP_GO_LIBS), -l$l)
WRAPPER_BIN_LIBS  = $(foreach l, $(WRAP_BIN_LIBS), -l$l)

RUN_TESTS = $(foreach b, $(SOL_TEST), printf "\n$(b)\n=====================\n" && ./$(b) &&) printf "==========\n"

lib: $(SONAME)

$(SONAME): $(SOL_SRC)
	cd $(LIB_DIR) &&\
		$(CXX) -c $(CXXFLAGS_LIB) $(SOL_SRC) $(WRAPPER_LIBS) -fPIC &&\
		ar -rcs libsolwrap.a *.o &&\
		rm *.o

lib-tests:
	$(CXX) $(CXXFLAGS) -o $(BUILD_DIR)/direct.test     $(SOLWRAP_DIR)/tests/direct_test.cpp         $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
	$(CXX) $(CXXFLAGS) -o $(BUILD_DIR)/persistent.test $(SOLWRAP_DIR)/tests/persistent_test.cpp     $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
	$(CXX) $(CXXFLAGS) -o $(BUILD_DIR)/subscribe.test  $(SOLWRAP_DIR)/tests/bulk_subscribe_test.cpp $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
# 	$(CC)  $(CXXFLAGS) -o $(BUILD_DIR)/c_client.test   $(SOLWRAP_DIR)/tests/c_client_test.c         $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
# 	$(CXX) $(CXXFLAGS) -o $(BUILD_DIR)/cache.test     $(SOLWRAP_DIR)/tests/cache_test.cpp           $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)

test:
	$(RUN_TESTS)

examples:
	CGO_LDFLAGS="$(STATIC_LIBS) $(WRAPPER_GO_LIBS)" CGO_CFLAGS="-fPIC $(CXXFLAGS)" go build -o pub solace.go gosol.go types.go publisher.queue.example.go
	CGO_LDFLAGS="$(STATIC_LIBS) $(WRAPPER_GO_LIBS)" CGO_CFLAGS="-fPIC $(CXXFLAGS)" go build -o sub solace.go gosol.go types.go subscriber.queue.example.go