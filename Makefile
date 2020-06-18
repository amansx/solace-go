LIBNAME        = solwrap
SONAME         = lib$(LIBNAME).so

SOLWRAP_DIR    = $(CURDIR)/solwrap
SOLCLIENT_DIR  = $(CURDIR)/solclient

BUILD_DIR      = $(CURDIR)/bin
TEST_DIR       = $(CURDIR)/test

PYINC=
PYLIB=
PYDEF=

INCDIRS              = $(SOLCLIENT_DIR)/include $(SOLWRAP_DIR)/include
LIBDIRS              = $(CURDIR)/lib $(SOLCLIENT_DIR)/lib
STATIC_LIBS          = $(CURDIR)/lib/libsolwrap.a $(SOLCLIENT_DIR)/lib/libsolclient.a
WRAP_LIBS            = pthread
WRAP_BIN_LIBS        = pthread rt
WRAP_GO_LIBS         = pthread rt stdc++
SYMS                 = PROVIDE_LOG_UTILITIES SOLCLIENT_CONST_PROPERTIES _REENTRANT _LINUX_X86_64
DEBUG                = -g

SOL_SRC     = $(wildcard $(CURDIR)/solwrap/src/*.cpp)
SOL_TEST    = $(wildcard ./bin/*.test)

# -Wl,--whole-archive libAlgatorc.a -Wl,--no-whole-archive
CXXFLAGS_LIB = $(foreach d, $(INCDIRS), -I$d) $(foreach s, $(SYMS), -D$s) -m64 $(DEBUG)
CXXFLAGS     = $(foreach d, $(INCDIRS), -I$d) $(foreach s, $(SYMS), -D$s) -m64 $(DEBUG)

WRAPPER_LIBS      = $(foreach l, $(WRAP_LIBS), -l$l)
WRAPPER_GO_LIBS   = $(foreach l, $(WRAP_GO_LIBS), -l$l)
WRAPPER_BIN_LIBS  = $(foreach l, $(WRAP_BIN_LIBS), -l$l)

RUN_TESTS    = $(foreach b, $(SOL_TEST), printf "\n$(b)\n=====================\n" && LD_LIBRARY_PATH=./bin/ ./$(b) ./src/solace.properties &&) printf "==========\n"

lib: $(SONAME)

$(SONAME): $(SOL_SRC)
	cd lib &&\
		$(CXX) -c $(CXXFLAGS_LIB) $(SOL_SRC) $(WRAPPER_LIBS) -fPIC &&\
		ar -rcs libsolwrap.a *.o &&\
		rm *.o

lib-tests:
# 	$(CC)  $(CXXFLAGS) -o $(TEST_DIR)/c_client.test   $(SOLWRAP_DIR)/tests/c_client_test.c         $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
# 	$(CXX) $(CXXFLAGS) -o $(TEST_DIR)/cache.test     $(SOLWRAP_DIR)/tests/cache_test.cpp           $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
	$(CXX) $(CXXFLAGS) -o $(TEST_DIR)/direct.test     $(SOLWRAP_DIR)/tests/direct_test.cpp         $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
	$(CXX) $(CXXFLAGS) -o $(TEST_DIR)/persistent.test $(SOLWRAP_DIR)/tests/persistent_test.cpp     $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)
	$(CXX) $(CXXFLAGS) -o $(TEST_DIR)/subscribe.test  $(SOLWRAP_DIR)/tests/bulk_subscribe_test.cpp $(STATIC_LIBS) $(WRAPPER_BIN_LIBS)

binding:
	GOPATH=$(CURDIR) CGO_LDFLAGS="$(STATIC_LIBS) $(WRAPPER_GO_LIBS)" CGO_CFLAGS="$(CXXFLAGS)" go install gosol
	GOPATH=$(CURDIR) CGO_LDFLAGS="$(STATIC_LIBS) $(WRAPPER_GO_LIBS)" CGO_CFLAGS="-fPIC $(CXXFLAGS)" go build -buildmode=plugin -o $(BUILD_DIR)/gosol.lib src/gosol.lib/main.go

binding-tests:
	GOPATH=$(CURDIR) go install gosol.test

test:
	$(RUN_TESTS)