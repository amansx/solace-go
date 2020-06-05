LIBNAME        = solwrap
SONAME         = lib$(LIBNAME).so

SOLWRAP_DIR    = $(CURDIR)/solwrap
SOLCLIENT_DIR  = $(CURDIR)/solclient

BUILD_DIR      = $(CURDIR)/bin

PYINC=
PYLIB=
PYDEF=

INCDIRS = $(SOLCLIENT_DIR)/include $(SOLWRAP_DIR)/include
LIBDIRS = $(SOLCLIENT_DIR)/lib $(SOLWRAP_DIR)/lib

SYMS    = PROVIDE_LOG_UTILITIES SOLCLIENT_CONST_PROPERTIES _REENTRANT _LINUX_X86_64
DEBUG   = -g

WRAP_LIBS  = solclient pthread
ALL_LIBS   = $(LIBNAME) solclient pthread

CXXFLAGS      = $(foreach d, $(INCDIRS), -I$d) $(foreach s, $(SYMS), -D$s) -m64 $(DEBUG)

LIBS_DIR     = $(foreach d, $(LIBDIRS), -L$d)
LIBS_WRAP    = $(foreach l, $(WRAP_LIBS), -l$l)
LIBS         = $(foreach l, $(ALL_LIBS), -l$l)

SOL_SRC     = $(wildcard solwrap/src/*.cpp)

lib: $(SONAME)

$(SONAME): $(SOL_SRC)
	$(CXX) -shared $(CXXFLAGS) $(LIBS_DIR) $(LIBS_WRAP) -fPIC -o $(SOLWRAP_DIR)/lib/$(SONAME) $(SOL_SRC)
	cp $(SOLWRAP_DIR)/lib/$(SONAME) $(BUILD_DIR)
	cp $(SOLCLIENT_DIR)/lib/libsolclient.so.1 $(BUILD_DIR)

lib-tests:
	$(CC)  $(CXXFLAGS) $(LIBS_DIR) -o $(BUILD_DIR)/test_c_client       $(SOLWRAP_DIR)/tests/c_client_test.c $(LIBS)
	$(CXX) $(CXXFLAGS) $(LIBS_DIR) -o $(BUILD_DIR)/test_direct         $(SOLWRAP_DIR)/tests/direct_test.cpp $(LIBS)
	$(CXX) $(CXXFLAGS) $(LIBS_DIR) -o $(BUILD_DIR)/test_persistent     $(SOLWRAP_DIR)/tests/persistent_test.cpp $(LIBS)
	$(CXX) $(CXXFLAGS) $(LIBS_DIR) -o $(BUILD_DIR)/test_cache          $(SOLWRAP_DIR)/tests/cache_test.cpp $(LIBS)
	$(CXX) $(CXXFLAGS) $(LIBS_DIR) -o $(BUILD_DIR)/test_subscribe      $(SOLWRAP_DIR)/tests/bulk_subscribe_test.cpp $(LIBS)

bindings:
	GOPATH=$(CURDIR) CGO_LDFLAGS="$(LIBS_DIR) $(LIBS)" CGO_CFLAGS="$(CXXFLAGS)" go install gosol

examples:
	GOPATH=$(CURDIR) CGO_LDFLAGS="$(LIBS_DIR) $(LIBS)" CGO_CFLAGS="$(CXXFLAGS)" go install solace-example