# Параметры прогресс-бара
PROGRESS_WIDTH := 50
PROGRESS_TITLE := "Building project"

define ProgressBarGoBuild
	awk -v bar_title=$(PROGRESS_TITLE) \
	    -v project=$(1) \
    	-v total_steps=$(2) \
    	-v bar_length=$(PROGRESS_WIDTH) \
	   	-f configs/_make_/lib/go-build/progress-bar.awk
endef

define GetTotals
	$(1) | grep -c '\<mkdir\>'
endef

define GoBuildProgress
	$(call CacheRAMDrive,$(PROJECT_NAME)/$(2))
	@steps=$$($(call GetTotals,$(call $(1),-n))) && \
	$(call $(1),) | \
	tee $(GO_BUILD_LOG) | \
	($(call ProgressBarGoBuild,$(2),$$steps) && \
	ERROR_COUNT=$$(grep -c '.go:' "$(GO_BUILD_LOG)"); \
	if [ "$$ERROR_COUNT" -gt 0 ]; then \
    	grep --color=always '.go:' $(GO_BUILD_LOG); \
    	rm -f $(GO_BUILD_LOG); \
    	exit 1; \
    fi)
endef