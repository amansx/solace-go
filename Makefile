windows:
	@make -f makefile.win
	@make -B lib-tests --no-print-directory -f makefile.win
	@make -B examples --no-print-directory -f makefile.win

linux:
	@make -f makefile.linux
	@make -B lib-tests --no-print-directory -f makefile.linux
	@make -B examples --no-print-directory -f makefile.linux

osx:
	make -f makefile.osx
	make -B lib-tests --no-print-directory -f makefile.osx
	make -B examples --no-print-directory -f makefile.osx