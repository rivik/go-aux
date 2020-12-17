projname := $(notdir $(CURDIR))

ver :=
buildno := 0
buildts := $(shell printf "%.3f" `date +%s.%N` | sed 's/\.//')

rev := $(notdir $(shell git describe --all --tags --long --exact-match --always --dirty --broken))
ifndef rev
$(warning Cannot detect git revision for build)
rev := unknown-revision
endif

adt_build_args_go :=
#adt_build_args_go := -a
adt_ldflags_go := -X github.com/rivik/go-aux/pkg/appver.LDBuildTS=$(buildts) -X github.com/rivik/go-aux/pkg/appver.LDBuildNo=$(buildno) -X github.com/rivik/go-aux/pkg/appver.LDSemVer=$(ver) -X github.com/rivik/go-aux/pkg/appver.LDAltVer=$(addsuffix -$(buildno),$(shell date +%Y.%m)) -X github.com/rivik/go-aux/pkg/appver.LDRevision=$(rev)

cmd_dir_go := cmd
tgt_prefix := .makefiletarget

projdir := $(CURDIR)
ifeq (,$(wildcard Makefile))
$(error Current directory $(shell pwd) is not a project root direcotry)
endif

assets := $(shell find assets/ -type f)
sources_go := $(shell find . -type f -name '*.go')
bin_dirs_go := $(filter %/, $(wildcard $(cmd_dir_go)/*/))
bin_names_go := $(notdir $(patsubst %/,%,$(bin_dirs_go)))
binaries_go := $(join $(bin_dirs_go),$(bin_names_go))

comma := ,
empty :=
space := $(empty) $(empty)
binaries_go_joined := $(subst $(space),$(comma),$(binaries_go))

.PHONY: all facts assets lint build clean
all: facts assets lint build
	@echo OK

facts:
	@echo Makefile running with config:
	@echo projname: $(projname)
	@echo projdir: $(projdir)
	@echo binaries_go: $(binaries_go)
	@echo binaries_go_joined: $(binaries_go_joined)
	@echo adt_build_args_go: $(adt_build_args_go)
	@echo adt_ldflags_go: $(adt_ldflags_go)
	@echo

assets: $(tgt_prefix).assets.$(rev)
$(tgt_prefix).assets.$(rev): $(assets)
	@echo ">> writing assets"
	# Un-setting GOOS and GOARCH here because the generated Go code is always the same,
	# but the cached object code is incompatible between architectures and OSes (which
	# breaks cross-building for different combinations on CI in the same container).
	if [ -d assets -a -f assets/assets_generate.go ]; then \
		cd assets && GOOS= GOARCH= go run assets_generate.go && cd .. ;\
	fi
	touch $@

build: $(binaries_go)
$(binaries_go): $(sources_go) $(assets)
	# for each binary from binaries_go: run build
	# with $@(cur_bin, like cmd/myprog/myprog)
	# and $(@D)(cur_bin_dir, like cmd/myprog)
	go build -o "$@" -v $(adt_build_args_go) -trimpath -tags netgo -ldflags '$(adt_ldflags_go) -w -s -extldflags "-static"' "./$(@D)"

lint: $(tgt_prefix).lint.$(rev)
$(tgt_prefix).lint.$(rev): $(sources_go)
	# https://github.com/golangci/golangci-lint#install
	golangci-lint run
	touch $@

clean:
	rm -f $(binaries_go)
	rm -f $(tgt_prefix).*
	rm -f go.sum
	go mod tidy -v
