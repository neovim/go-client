package main

type API struct {
	BufAddDecoration    *Type `json:"nvim__buf_add_decoration" msgpack:"nvim__buf_add_decoration"`
	BufRedrawRange      *Type `json:"nvim__buf_redraw_range" msgpack:"nvim__buf_redraw_range"`
	BufSetLuahl         *Type `json:"nvim__buf_set_luahl" msgpack:"nvim__buf_set_luahl"`
	BufStats            *Type `json:"nvim__buf_stats" msgpack:"nvim__buf_stats"`
	GetLibDir           *Type `json:"nvim__get_lib_dir" msgpack:"nvim__get_lib_dir"`
	ID                  *Type `json:"nvim__id" msgpack:"nvim__id"`
	IDArray             *Type `json:"nvim__id_array" msgpack:"nvim__id_array"`
	IDDictionary        *Type `json:"nvim__id_dictionary" msgpack:"nvim__id_dictionary"`
	IDFloat             *Type `json:"nvim__id_float" msgpack:"nvim__id_float"`
	InspectCell         *Type `json:"nvim__inspect_cell" msgpack:"nvim__inspect_cell"`
	PutAttr             *Type `json:"nvim__put_attr" msgpack:"nvim__put_attr"`
	Stats               *Type `json:"nvim__stats" msgpack:"nvim__stats"`
	BufAddHighlight     *Type `json:"nvim_buf_add_highlight" msgpack:"nvim_buf_add_highlight"`
	BufAttach           *Type `json:"nvim_buf_attach" msgpack:"nvim_buf_attach"`
	BufClearNamespace   *Type `json:"nvim_buf_clear_namespace" msgpack:"nvim_buf_clear_namespace"`
	BufDelExtmark       *Type `json:"nvim_buf_del_extmark" msgpack:"nvim_buf_del_extmark"`
	BufDelKeymap        *Type `json:"nvim_buf_del_keymap" msgpack:"nvim_buf_del_keymap"`
	BufDelVar           *Type `json:"nvim_buf_del_var" msgpack:"nvim_buf_del_var"`
	BufDetach           *Type `json:"nvim_buf_detach" msgpack:"nvim_buf_detach"`
	BufGetChangedtick   *Type `json:"nvim_buf_get_changedtick" msgpack:"nvim_buf_get_changedtick"`
	BufGetCommands      *Type `json:"nvim_buf_get_commands" msgpack:"nvim_buf_get_commands"`
	BufGetExtmarkByID   *Type `json:"nvim_buf_get_extmark_by_id" msgpack:"nvim_buf_get_extmark_by_id"`
	BufGetExtmarks      *Type `json:"nvim_buf_get_extmarks" msgpack:"nvim_buf_get_extmarks"`
	BufGetKeymap        *Type `json:"nvim_buf_get_keymap" msgpack:"nvim_buf_get_keymap"`
	BufGetLines         *Type `json:"nvim_buf_get_lines" msgpack:"nvim_buf_get_lines"`
	BufGetMark          *Type `json:"nvim_buf_get_mark" msgpack:"nvim_buf_get_mark"`
	BufGetName          *Type `json:"nvim_buf_get_name" msgpack:"nvim_buf_get_name"`
	BufGetOffset        *Type `json:"nvim_buf_get_offset" msgpack:"nvim_buf_get_offset"`
	BufGetOption        *Type `json:"nvim_buf_get_option" msgpack:"nvim_buf_get_option"`
	BufGetVar           *Type `json:"nvim_buf_get_var" msgpack:"nvim_buf_get_var"`
	BufGetVirtualText   *Type `json:"nvim_buf_get_virtual_text" msgpack:"nvim_buf_get_virtual_text"`
	BufIsLoaded         *Type `json:"nvim_buf_is_loaded" msgpack:"nvim_buf_is_loaded"`
	BufIsValid          *Type `json:"nvim_buf_is_valid" msgpack:"nvim_buf_is_valid"`
	BufLineCount        *Type `json:"nvim_buf_line_count" msgpack:"nvim_buf_line_count"`
	BufSetExtmark       *Type `json:"nvim_buf_set_extmark" msgpack:"nvim_buf_set_extmark"`
	BufSetKeymap        *Type `json:"nvim_buf_set_keymap" msgpack:"nvim_buf_set_keymap"`
	BufSetLines         *Type `json:"nvim_buf_set_lines" msgpack:"nvim_buf_set_lines"`
	BufSetName          *Type `json:"nvim_buf_set_name" msgpack:"nvim_buf_set_name"`
	BufSetOption        *Type `json:"nvim_buf_set_option" msgpack:"nvim_buf_set_option"`
	BufSetVar           *Type `json:"nvim_buf_set_var" msgpack:"nvim_buf_set_var"`
	BufSetVirtualText   *Type `json:"nvim_buf_set_virtual_text" msgpack:"nvim_buf_set_virtual_text"`
	CallAtomic          *Type `json:"nvim_call_atomic" msgpack:"nvim_call_atomic"`
	CallDictFunction    *Type `json:"nvim_call_dict_function" msgpack:"nvim_call_dict_function"`
	CallFunction        *Type `json:"nvim_call_function" msgpack:"nvim_call_function"`
	Command             *Type `json:"nvim_command" msgpack:"nvim_command"`
	CreateBuf           *Type `json:"nvim_create_buf" msgpack:"nvim_create_buf"`
	CreateNamespace     *Type `json:"nvim_create_namespace" msgpack:"nvim_create_namespace"`
	DelCurrentLine      *Type `json:"nvim_del_current_line" msgpack:"nvim_del_current_line"`
	DelKeymap           *Type `json:"nvim_del_keymap" msgpack:"nvim_del_keymap"`
	DelVar              *Type `json:"nvim_del_var" msgpack:"nvim_del_var"`
	ErrWrite            *Type `json:"nvim_err_write" msgpack:"nvim_err_write"`
	ErrWriteln          *Type `json:"nvim_err_writeln" msgpack:"nvim_err_writeln"`
	Eval                *Type `json:"nvim_eval" msgpack:"nvim_eval"`
	Exec                *Type `json:"nvim_exec" msgpack:"nvim_exec"`
	ExecLua             *Type `json:"nvim_exec_lua" msgpack:"nvim_exec_lua"`
	Feedkeys            *Type `json:"nvim_feedkeys" msgpack:"nvim_feedkeys"`
	GetAPIInfo          *Type `json:"nvim_get_api_info" msgpack:"nvim_get_api_info"`
	GetChanInfo         *Type `json:"nvim_get_chan_info" msgpack:"nvim_get_chan_info"`
	GetColorByName      *Type `json:"nvim_get_color_by_name" msgpack:"nvim_get_color_by_name"`
	GetColorMap         *Type `json:"nvim_get_color_map" msgpack:"nvim_get_color_map"`
	GetCommands         *Type `json:"nvim_get_commands" msgpack:"nvim_get_commands"`
	GetContext          *Type `json:"nvim_get_context" msgpack:"nvim_get_context"`
	GetCurrentBuf       *Type `json:"nvim_get_current_buf" msgpack:"nvim_get_current_buf"`
	GetCurrentLine      *Type `json:"nvim_get_current_line" msgpack:"nvim_get_current_line"`
	GetCurrentTabpage   *Type `json:"nvim_get_current_tabpage" msgpack:"nvim_get_current_tabpage"`
	GetCurrentWin       *Type `json:"nvim_get_current_win" msgpack:"nvim_get_current_win"`
	GetHlByID           *Type `json:"nvim_get_hl_by_id" msgpack:"nvim_get_hl_by_id"`
	GetHlByName         *Type `json:"nvim_get_hl_by_name" msgpack:"nvim_get_hl_by_name"`
	GetHlIDByName       *Type `json:"nvim_get_hl_id_by_name" msgpack:"nvim_get_hl_id_by_name"`
	GetKeymap           *Type `json:"nvim_get_keymap" msgpack:"nvim_get_keymap"`
	GetMode             *Type `json:"nvim_get_mode" msgpack:"nvim_get_mode"`
	GetNamespaces       *Type `json:"nvim_get_namespaces" msgpack:"nvim_get_namespaces"`
	GetOption           *Type `json:"nvim_get_option" msgpack:"nvim_get_option"`
	GetProc             *Type `json:"nvim_get_proc" msgpack:"nvim_get_proc"`
	GetProcChildren     *Type `json:"nvim_get_proc_children" msgpack:"nvim_get_proc_children"`
	GetRuntimeFile      *Type `json:"nvim_get_runtime_file" msgpack:"nvim_get_runtime_file"`
	GetVar              *Type `json:"nvim_get_var" msgpack:"nvim_get_var"`
	GetVvar             *Type `json:"nvim_get_vvar" msgpack:"nvim_get_vvar"`
	Input               *Type `json:"nvim_input" msgpack:"nvim_input"`
	InputMouse          *Type `json:"nvim_input_mouse" msgpack:"nvim_input_mouse"`
	ListBufs            *Type `json:"nvim_list_bufs" msgpack:"nvim_list_bufs"`
	ListChans           *Type `json:"nvim_list_chans" msgpack:"nvim_list_chans"`
	ListRuntimePaths    *Type `json:"nvim_list_runtime_paths" msgpack:"nvim_list_runtime_paths"`
	ListTabpages        *Type `json:"nvim_list_tabpages" msgpack:"nvim_list_tabpages"`
	ListUis             *Type `json:"nvim_list_uis" msgpack:"nvim_list_uis"`
	ListWins            *Type `json:"nvim_list_wins" msgpack:"nvim_list_wins"`
	LoadContext         *Type `json:"nvim_load_context" msgpack:"nvim_load_context"`
	OpenWin             *Type `json:"nvim_open_win" msgpack:"nvim_open_win"`
	OutWrite            *Type `json:"nvim_out_write" msgpack:"nvim_out_write"`
	ParseExpression     *Type `json:"nvim_parse_expression" msgpack:"nvim_parse_expression"`
	Paste               *Type `json:"nvim_paste" msgpack:"nvim_paste"`
	Put                 *Type `json:"nvim_put" msgpack:"nvim_put"`
	ReplaceTermcodes    *Type `json:"nvim_replace_termcodes" msgpack:"nvim_replace_termcodes"`
	SelectPopupmenuItem *Type `json:"nvim_select_popupmenu_item" msgpack:"nvim_select_popupmenu_item"`
	SetClientInfo       *Type `json:"nvim_set_client_info" msgpack:"nvim_set_client_info"`
	SetCurrentBuf       *Type `json:"nvim_set_current_buf" msgpack:"nvim_set_current_buf"`
	SetCurrentDir       *Type `json:"nvim_set_current_dir" msgpack:"nvim_set_current_dir"`
	SetCurrentLine      *Type `json:"nvim_set_current_line" msgpack:"nvim_set_current_line"`
	SetCurrentTabpage   *Type `json:"nvim_set_current_tabpage" msgpack:"nvim_set_current_tabpage"`
	SetCurrentWin       *Type `json:"nvim_set_current_win" msgpack:"nvim_set_current_win"`
	SetKeymap           *Type `json:"nvim_set_keymap" msgpack:"nvim_set_keymap"`
	SetOption           *Type `json:"nvim_set_option" msgpack:"nvim_set_option"`
	SetVar              *Type `json:"nvim_set_var" msgpack:"nvim_set_var"`
	SetVvar             *Type `json:"nvim_set_vvar" msgpack:"nvim_set_vvar"`
	Strwidth            *Type `json:"nvim_strwidth" msgpack:"nvim_strwidth"`
	Subscribe           *Type `json:"nvim_subscribe" msgpack:"nvim_subscribe"`
	TabpageDelVar       *Type `json:"nvim_tabpage_del_var" msgpack:"nvim_tabpage_del_var"`
	TabpageGetNumber    *Type `json:"nvim_tabpage_get_number" msgpack:"nvim_tabpage_get_number"`
	TabpageGetVar       *Type `json:"nvim_tabpage_get_var" msgpack:"nvim_tabpage_get_var"`
	TabpageGetWin       *Type `json:"nvim_tabpage_get_win" msgpack:"nvim_tabpage_get_win"`
	TabpageIsValid      *Type `json:"nvim_tabpage_is_valid" msgpack:"nvim_tabpage_is_valid"`
	TabpageListWins     *Type `json:"nvim_tabpage_list_wins" msgpack:"nvim_tabpage_list_wins"`
	TabpageSetVar       *Type `json:"nvim_tabpage_set_var" msgpack:"nvim_tabpage_set_var"`
	UIAttach            *Type `json:"nvim_ui_attach" msgpack:"nvim_ui_attach"`
	UIDetach            *Type `json:"nvim_ui_detach" msgpack:"nvim_ui_detach"`
	UIPumSetHeight      *Type `json:"nvim_ui_pum_set_height" msgpack:"nvim_ui_pum_set_height"`
	UISetOption         *Type `json:"nvim_ui_set_option" msgpack:"nvim_ui_set_option"`
	UITryResize         *Type `json:"nvim_ui_try_resize" msgpack:"nvim_ui_try_resize"`
	UITryResizeGrid     *Type `json:"nvim_ui_try_resize_grid" msgpack:"nvim_ui_try_resize_grid"`
	Unsubscribe         *Type `json:"nvim_unsubscribe" msgpack:"nvim_unsubscribe"`
	WinClose            *Type `json:"nvim_win_close" msgpack:"nvim_win_close"`
	WinDelVar           *Type `json:"nvim_win_del_var" msgpack:"nvim_win_del_var"`
	WinGetBuf           *Type `json:"nvim_win_get_buf" msgpack:"nvim_win_get_buf"`
	WinGetConfig        *Type `json:"nvim_win_get_config" msgpack:"nvim_win_get_config"`
	WinGetCursor        *Type `json:"nvim_win_get_cursor" msgpack:"nvim_win_get_cursor"`
	WinGetHeight        *Type `json:"nvim_win_get_height" msgpack:"nvim_win_get_height"`
	WinGetNumber        *Type `json:"nvim_win_get_number" msgpack:"nvim_win_get_number"`
	WinGetOption        *Type `json:"nvim_win_get_option" msgpack:"nvim_win_get_option"`
	WinGetPosition      *Type `json:"nvim_win_get_position" msgpack:"nvim_win_get_position"`
	WinGetTabpage       *Type `json:"nvim_win_get_tabpage" msgpack:"nvim_win_get_tabpage"`
	WinGetVar           *Type `json:"nvim_win_get_var" msgpack:"nvim_win_get_var"`
	WinGetWidth         *Type `json:"nvim_win_get_width" msgpack:"nvim_win_get_width"`
	WinIsValid          *Type `json:"nvim_win_is_valid" msgpack:"nvim_win_is_valid"`
	WinSetBuf           *Type `json:"nvim_win_set_buf" msgpack:"nvim_win_set_buf"`
	WinSetConfig        *Type `json:"nvim_win_set_config" msgpack:"nvim_win_set_config"`
	WinSetCursor        *Type `json:"nvim_win_set_cursor" msgpack:"nvim_win_set_cursor"`
	WinSetHeight        *Type `json:"nvim_win_set_height" msgpack:"nvim_win_set_height"`
	WinSetOption        *Type `json:"nvim_win_set_option" msgpack:"nvim_win_set_option"`
	WinSetVar           *Type `json:"nvim_win_set_var" msgpack:"nvim_win_set_var"`
	WinSetWidth         *Type `json:"nvim_win_set_width" msgpack:"nvim_win_set_width"`
}

type Type struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           []interface{} `json:"doc,omitempty" msgpack:"doc,omitempty"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl,omitempty" msgpack:"c_decl,omitempty"`
}

type BufAddDecoration struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type BufRedrawRange struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type BufSetLuahl struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type BufStats struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type GetLibDir struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type ID struct {
	Annotations   []interface{}   `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}     `json:"doc" msgpack:"doc"`
	Parameters    []interface{}   `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc IDParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string        `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}   `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string          `json:"signature" msgpack:"signature"`
	Declaration   string          `json:"c_decl" msgpack:"c_decl"`
}

type IDParametersDoc struct {
	Obj string `json:"obj" msgpack:"obj"`
}

type IDArray struct {
	Annotations   []interface{}        `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}          `json:"doc" msgpack:"doc"`
	Parameters    []interface{}        `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc IDArrayParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}        `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string               `json:"signature" msgpack:"signature"`
	Declaration   string               `json:"c_decl" msgpack:"c_decl"`
}

type IDArrayParametersDoc struct {
	Arr string `json:"arr" msgpack:"arr"`
}

type IDDictionary struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc IDDictionaryParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
	Declaration   string                    `json:"c_decl" msgpack:"c_decl"`
}

type IDDictionaryParametersDoc struct {
	Dct string `json:"dct" msgpack:"dct"`
}

type IDFloat struct {
	Annotations   []interface{}        `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}          `json:"doc" msgpack:"doc"`
	Parameters    []interface{}        `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc IDFloatParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}        `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string               `json:"signature" msgpack:"signature"`
	Declaration   string               `json:"c_decl" msgpack:"c_decl"`
}

type IDFloatParametersDoc struct {
	Flt string `json:"flt" msgpack:"flt"`
}

type InspectCell struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type PutAttr struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type Stats struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
	Declaration   string        `json:"c_decl" msgpack:"c_decl"`
}

type BufAddHighlight struct {
	Annotations   []interface{}                `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                  `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufAddHighlightParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                     `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                       `json:"signature" msgpack:"signature"`
	Declaration   string                       `json:"c_decl" msgpack:"c_decl"`
}

type BufAddHighlightParametersDoc struct {
	Buffer   string `json:"buffer" msgpack:"buffer"`
	ColEnd   string `json:"col_end" msgpack:"col_end"`
	ColStart string `json:"col_start" msgpack:"col_start"`
	HlGroup  string `json:"hl_group" msgpack:"hl_group"`
	Line     string `json:"line" msgpack:"line"`
	NsID     string `json:"ns_id" msgpack:"ns_id"`
}

type BufAttach struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufAttachParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
	Declaration   string                 `json:"c_decl" msgpack:"c_decl"`
}

type BufAttachParametersDoc struct {
	Buffer     string `json:"buffer" msgpack:"buffer"`
	Opts       string `json:"opts" msgpack:"opts"`
	SendBuffer string `json:"send_buffer" msgpack:"send_buffer"`
}

type BufClearNamespace struct {
	Annotations   []interface{}                  `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                    `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                  `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufClearNamespaceParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                         `json:"signature" msgpack:"signature"`
	Declaration   string                         `json:"c_decl" msgpack:"c_decl"`
}

type BufClearNamespaceParametersDoc struct {
	Buffer    string `json:"buffer" msgpack:"buffer"`
	LineEnd   string `json:"line_end" msgpack:"line_end"`
	LineStart string `json:"line_start" msgpack:"line_start"`
	NsID      string `json:"ns_id" msgpack:"ns_id"`
}

type BufDelExtmark struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufDelExtmarkParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                   `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
	Declaration   string                     `json:"c_decl" msgpack:"c_decl"`
}

type BufDelExtmarkParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	ID     string `json:"id" msgpack:"id"`
	NsID   string `json:"ns_id" msgpack:"ns_id"`
}

type BufDelKeymap struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufDelKeymapParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
	Declaration   string                    `json:"c_decl" msgpack:"c_decl"`
}

type BufDelKeymapParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufDelVar struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufDelVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
	Declaration   string                 `json:"c_decl" msgpack:"c_decl"`
}

type BufDelVarParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Name   string `json:"name" msgpack:"name"`
}

type BufDetach struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufDetachParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
	Declaration   string                 `json:"c_decl" msgpack:"c_decl"`
}

type BufDetachParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufGetChangedtick struct {
	Annotations   []interface{}                  `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                    `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                  `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetChangedtickParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                       `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                         `json:"signature" msgpack:"signature"`
	Declaration   string                         `json:"c_decl" msgpack:"c_decl"`
}

type BufGetChangedtickParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufGetCommands struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetCommandsParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
	Declaration   string                      `json:"c_decl" msgpack:"c_decl"`
}

type BufGetCommandsParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Opts   string `json:"opts" msgpack:"opts"`
}

type BufGetExtmarkByID struct {
	Annotations   []interface{}                  `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                    `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                  `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetExtmarkByIDParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                       `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                         `json:"signature" msgpack:"signature"`
	Declaration   string                         `json:"c_decl" msgpack:"c_decl"`
}

type BufGetExtmarkByIDParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	ID     string `json:"id" msgpack:"id"`
	NsID   string `json:"ns_id" msgpack:"ns_id"`
}

type BufGetExtmarks struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetExtmarksParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
	Declaration   string                      `json:"c_decl" msgpack:"c_decl"`
}

type BufGetExtmarksParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	End    string `json:"end" msgpack:"end"`
	NsID   string `json:"ns_id" msgpack:"ns_id"`
	Opts   string `json:"opts" msgpack:"opts"`
	Start  string `json:"start" msgpack:"start"`
}

type BufGetKeymap struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetKeymapParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
	Declaration   string                    `json:"c_decl" msgpack:"c_decl"`
}

type BufGetKeymapParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Mode   string `json:"mode" msgpack:"mode"`
}

type BufGetLines struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetLinesParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                 `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
	Declaration   string                   `json:"c_decl" msgpack:"c_decl"`
}

type BufGetLinesParametersDoc struct {
	Buffer         string `json:"buffer" msgpack:"buffer"`
	End            string `json:"end" msgpack:"end"`
	Start          string `json:"start" msgpack:"start"`
	StrictIndexing string `json:"strict_indexing" msgpack:"strict_indexing"`
}

type BufGetMark struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetMarkParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
	Declaration   string                  `json:"c_decl" msgpack:"c_decl"`
}

type BufGetMarkParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Name   string `json:"name" msgpack:"name"`
}

type BufGetName struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetNameParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
	Declaration   string                  `json:"c_decl" msgpack:"c_decl"`
}

type BufGetNameParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufGetOffset struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetOffsetParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
	Declaration   string                    `json:"c_decl" msgpack:"c_decl"`
}

type BufGetOffsetParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Index  string `json:"index" msgpack:"index"`
}

type BufGetOption struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetOptionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
	Declaration   string                    `json:"c_decl" msgpack:"c_decl"`
}

type BufGetOptionParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Name   string `json:"name" msgpack:"name"`
}

type BufGetVar struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
	Declaration   string                 `json:"c_decl" msgpack:"c_decl"`
}

type BufGetVarParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Name   string `json:"name" msgpack:"name"`
}

type BufGetVirtualText struct {
	Annotations   []interface{}                  `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                    `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                  `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufGetVirtualTextParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                       `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                         `json:"signature" msgpack:"signature"`
}

type BufGetVirtualTextParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Line   string `json:"line" msgpack:"line"`
}

type BufIsLoaded struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufIsLoadedParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                 `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type BufIsLoadedParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufIsValid struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufIsValidParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
}

type BufIsValidParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufLineCount struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufLineCountParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type BufLineCountParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufSetExtmark struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufSetExtmarkParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                   `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type BufSetExtmarkParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Col    string `json:"col" msgpack:"col"`
	ID     string `json:"id" msgpack:"id"`
	Line   string `json:"line" msgpack:"line"`
	NsID   string `json:"ns_id" msgpack:"ns_id"`
	Opts   string `json:"opts" msgpack:"opts"`
}

type BufSetKeymap struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufSetKeymapParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type BufSetKeymapParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type BufSetLines struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufSetLinesParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}            `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type BufSetLinesParametersDoc struct {
	Buffer         string `json:"buffer" msgpack:"buffer"`
	End            string `json:"end" msgpack:"end"`
	Replacement    string `json:"replacement" msgpack:"replacement"`
	Start          string `json:"start" msgpack:"start"`
	StrictIndexing string `json:"strict_indexing" msgpack:"strict_indexing"`
}

type BufSetName struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufSetNameParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}           `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
}

type BufSetNameParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Name   string `json:"name" msgpack:"name"`
}

type BufSetOption struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufSetOptionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type BufSetOptionParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Name   string `json:"name" msgpack:"name"`
	Value  string `json:"value" msgpack:"value"`
}

type BufSetVar struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufSetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type BufSetVarParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Name   string `json:"name" msgpack:"name"`
	Value  string `json:"value" msgpack:"value"`
}

type BufSetVirtualText struct {
	Annotations   []interface{}                  `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                    `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                  `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc BufSetVirtualTextParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                       `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                         `json:"signature" msgpack:"signature"`
}

type BufSetVirtualTextParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Chunks string `json:"chunks" msgpack:"chunks"`
	Line   string `json:"line" msgpack:"line"`
	NsID   string `json:"ns_id" msgpack:"ns_id"`
	Opts   string `json:"opts" msgpack:"opts"`
}

type CallAtomic struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc CallAtomicParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
}

type CallAtomicParametersDoc struct {
	Calls string `json:"calls" msgpack:"calls"`
}

type CallDictFunction struct {
	Annotations   []interface{}                 `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                   `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                 `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc CallDictFunctionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                 `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                        `json:"signature" msgpack:"signature"`
}

type CallDictFunctionParametersDoc struct {
	Args string `json:"args" msgpack:"args"`
	Dict string `json:"dict" msgpack:"dict"`
	Fn   string `json:"fn" msgpack:"fn"`
}

type CallFunction struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc CallFunctionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type CallFunctionParametersDoc struct {
	Args string `json:"args" msgpack:"args"`
	Fn   string `json:"fn" msgpack:"fn"`
}

type Command struct {
	Annotations   []interface{}        `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}          `json:"doc" msgpack:"doc"`
	Parameters    []interface{}        `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc CommandParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}        `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string               `json:"signature" msgpack:"signature"`
}

type CommandParametersDoc struct {
	Command string `json:"command" msgpack:"command"`
}

type CreateBuf struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc CreateBufParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type CreateBufParametersDoc struct {
	Listed  string `json:"listed" msgpack:"listed"`
	Scratch string `json:"scratch" msgpack:"scratch"`
}

type CreateNamespace struct {
	Annotations   []interface{}                `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                  `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc CreateNamespaceParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                     `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                       `json:"signature" msgpack:"signature"`
}

type CreateNamespaceParametersDoc struct {
	Name string `json:"name" msgpack:"name"`
}

type DelCurrentLine struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type DelKeymap struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string      `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type DelVar struct {
	Annotations   []interface{}       `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}         `json:"doc" msgpack:"doc"`
	Parameters    []interface{}       `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc DelVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}       `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}       `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string              `json:"signature" msgpack:"signature"`
}

type DelVarParametersDoc struct {
	Name string `json:"name" msgpack:"name"`
}

type ErrWrite struct {
	Annotations   []interface{}         `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}           `json:"doc" msgpack:"doc"`
	Parameters    []interface{}         `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc ErrWriteParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}         `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}         `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                `json:"signature" msgpack:"signature"`
}

type ErrWriteParametersDoc struct {
	Str string `json:"str" msgpack:"str"`
}

type ErrWriteln struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc ErrWritelnParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}           `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string                `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
}

type ErrWritelnParametersDoc struct {
	Str string `json:"str" msgpack:"str"`
}

type Eval struct {
	Annotations   []interface{}     `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}       `json:"doc" msgpack:"doc"`
	Parameters    []interface{}     `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc EvalParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}     `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string            `json:"signature" msgpack:"signature"`
}

type EvalParametersDoc struct {
	Expr string `json:"expr" msgpack:"expr"`
}

type Exec struct {
	Annotations   []interface{}     `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}       `json:"doc" msgpack:"doc"`
	Parameters    []interface{}     `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc ExecParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string            `json:"signature" msgpack:"signature"`
}

type ExecParametersDoc struct {
	Output string `json:"output" msgpack:"output"`
	Src    string `json:"src" msgpack:"src"`
}

type ExecLua struct {
	Annotations   []interface{}        `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}          `json:"doc" msgpack:"doc"`
	Parameters    []interface{}        `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc ExecLuaParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}        `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string               `json:"signature" msgpack:"signature"`
}

type ExecLuaParametersDoc struct {
	Args string `json:"args" msgpack:"args"`
	Code string `json:"code" msgpack:"code"`
}

type Feedkeys struct {
	Annotations   []interface{}         `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}           `json:"doc" msgpack:"doc"`
	Parameters    []interface{}         `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc FeedkeysParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}         `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                `json:"signature" msgpack:"signature"`
}

type FeedkeysParametersDoc struct {
	EscapeCsi string `json:"escape_csi" msgpack:"escape_csi"`
	Keys      string `json:"keys" msgpack:"keys"`
	Mode      string `json:"mode" msgpack:"mode"`
}

type GetAPIInfo struct {
	Annotations   []string      `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetChanInfo struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetColorByName struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetColorByNameParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
}

type GetColorByNameParametersDoc struct {
	Name string `json:"name" msgpack:"name"`
}

type GetColorMap struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetColorMapParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                 `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type GetColorMapParametersDoc struct {
}

type GetCommands struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetCommandsParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                 `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type GetCommandsParametersDoc struct {
	Opts string `json:"opts" msgpack:"opts"`
}

type GetContext struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetContextParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
}

type GetContextParametersDoc struct {
	Opts string `json:"opts" msgpack:"opts"`
}

type GetCurrentBuf struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetCurrentLine struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetCurrentTabpage struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetCurrentWin struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetHlByID struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetHlByIDParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type GetHlByIDParametersDoc struct {
	HlID string `json:"hl_id" msgpack:"hl_id"`
	Rgb  string `json:"rgb" msgpack:"rgb"`
}

type GetHlByName struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetHlByNameParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                 `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string                 `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type GetHlByNameParametersDoc struct {
	Name string `json:"name" msgpack:"name"`
	Rgb  string `json:"rgb" msgpack:"rgb"`
}

type GetHlIDByName struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetKeymap struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetKeymapParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type GetKeymapParametersDoc struct {
	Mode string `json:"mode" msgpack:"mode"`
}

type GetMode struct {
	Annotations   []string      `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetNamespaces struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetOption struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetOptionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type GetOptionParametersDoc struct {
	Name string `json:"name" msgpack:"name"`
}

type GetProc struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetProcChildren struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type GetRuntimeFile struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetRuntimeFileParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
}

type GetRuntimeFileParametersDoc struct {
	All  string `json:"all" msgpack:"all"`
	Name string `json:"name" msgpack:"name"`
}

type GetVar struct {
	Annotations   []interface{}       `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}         `json:"doc" msgpack:"doc"`
	Parameters    []interface{}       `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string            `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}       `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string              `json:"signature" msgpack:"signature"`
}

type GetVarParametersDoc struct {
	Name string `json:"name" msgpack:"name"`
}

type GetVvar struct {
	Annotations   []interface{}        `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}          `json:"doc" msgpack:"doc"`
	Parameters    []interface{}        `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc GetVvarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}        `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string               `json:"signature" msgpack:"signature"`
}

type GetVvarParametersDoc struct {
	Name string `json:"name" msgpack:"name"`
}

type Input struct {
	Annotations   []string           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}        `json:"doc" msgpack:"doc"`
	Parameters    []interface{}      `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc InputParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string           `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}      `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string             `json:"signature" msgpack:"signature"`
}

type InputParametersDoc struct {
	Keys string `json:"keys" msgpack:"keys"`
}

type InputMouse struct {
	Annotations   []string                `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc InputMouseParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}           `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
}

type InputMouseParametersDoc struct {
	Action   string `json:"action" msgpack:"action"`
	Button   string `json:"button" msgpack:"button"`
	Col      string `json:"col" msgpack:"col"`
	Grid     string `json:"grid" msgpack:"grid"`
	Modifier string `json:"modifier" msgpack:"modifier"`
	Row      string `json:"row" msgpack:"row"`
}

type ListBufs struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type ListChans struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type ListRuntimePaths struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type ListTabpages struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type ListUis struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type ListWins struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []string      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type LoadContext struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc LoadContextParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}            `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type LoadContextParametersDoc struct {
	Dict string `json:"dict" msgpack:"dict"`
}

type OpenWin struct {
	Annotations   []interface{}        `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}          `json:"doc" msgpack:"doc"`
	Parameters    []interface{}        `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc OpenWinParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}        `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string               `json:"signature" msgpack:"signature"`
}

type OpenWinParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Config string `json:"config" msgpack:"config"`
	Enter  string `json:"enter" msgpack:"enter"`
}

type OutWrite struct {
	Annotations   []interface{}         `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}           `json:"doc" msgpack:"doc"`
	Parameters    []interface{}         `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc OutWriteParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}         `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}         `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                `json:"signature" msgpack:"signature"`
}

type OutWriteParametersDoc struct {
	Str string `json:"str" msgpack:"str"`
}

type ParseExpression struct {
	Annotations   []string                     `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                  `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc ParseExpressionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                     `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                       `json:"signature" msgpack:"signature"`
}

type ParseExpressionParametersDoc struct {
	Expr      string `json:"expr" msgpack:"expr"`
	Flags     string `json:"flags" msgpack:"flags"`
	Highlight string `json:"highlight" msgpack:"highlight"`
}

type Paste struct {
	Annotations   []interface{}      `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}        `json:"doc" msgpack:"doc"`
	Parameters    []interface{}      `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc PasteParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string           `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}      `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string             `json:"signature" msgpack:"signature"`
}

type PasteParametersDoc struct {
	Crlf  string `json:"crlf" msgpack:"crlf"`
	Data  string `json:"data" msgpack:"data"`
	Phase string `json:"phase" msgpack:"phase"`
}

type Put struct {
	Annotations   []interface{}    `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}      `json:"doc" msgpack:"doc"`
	Parameters    []interface{}    `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc PutParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}    `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string           `json:"signature" msgpack:"signature"`
}

type PutParametersDoc struct {
	After  string `json:"after" msgpack:"after"`
	Follow string `json:"follow" msgpack:"follow"`
	Lines  string `json:"lines" msgpack:"lines"`
	Type   string `json:"type" msgpack:"type"`
}

type ReplaceTermcodes struct {
	Annotations   []interface{}                 `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                   `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                 `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc ReplaceTermcodesParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}                 `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string                      `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                        `json:"signature" msgpack:"signature"`
}

type ReplaceTermcodesParametersDoc struct {
	DoLt     string `json:"do_lt" msgpack:"do_lt"`
	FromPart string `json:"from_part" msgpack:"from_part"`
	Special  string `json:"special" msgpack:"special"`
	Str      string `json:"str" msgpack:"str"`
}

type SelectPopupmenuItem struct {
	Annotations   []interface{}                    `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                      `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                    `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SelectPopupmenuItemParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}                    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                    `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                           `json:"signature" msgpack:"signature"`
}

type SelectPopupmenuItemParametersDoc struct {
	Finish string `json:"finish" msgpack:"finish"`
	Insert string `json:"insert" msgpack:"insert"`
	Item   string `json:"item" msgpack:"item"`
	Opts   string `json:"opts" msgpack:"opts"`
}

type SetClientInfo struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetClientInfoParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}              `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type SetClientInfoParametersDoc struct {
	Attributes string `json:"attributes" msgpack:"attributes"`
	Methods    string `json:"methods" msgpack:"methods"`
	Name       string `json:"name" msgpack:"name"`
	Type       string `json:"type" msgpack:"type"`
	Version    string `json:"version" msgpack:"version"`
}

type SetCurrentBuf struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetCurrentBufParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}              `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type SetCurrentBufParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
}

type SetCurrentDir struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetCurrentDirParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}              `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type SetCurrentDirParametersDoc struct {
	Dir string `json:"dir" msgpack:"dir"`
}

type SetCurrentLine struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetCurrentLineParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
}

type SetCurrentLineParametersDoc struct {
	Line string `json:"line" msgpack:"line"`
}

type SetCurrentTabpage struct {
	Annotations   []interface{}                  `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                    `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                  `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetCurrentTabpageParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                         `json:"signature" msgpack:"signature"`
}

type SetCurrentTabpageParametersDoc struct {
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
}

type SetCurrentWin struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetCurrentWinParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}              `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type SetCurrentWinParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type SetKeymap struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetKeymapParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type SetKeymapParametersDoc struct {
	LHS  string `json:"lhs" msgpack:"lhs"`
	Mode string `json:"mode" msgpack:"mode"`
	Opts string `json:"opts" msgpack:"opts"`
	RHS  string `json:"rhs" msgpack:"rhs"`
}

type SetOption struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetOptionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type SetOptionParametersDoc struct {
	Name  string `json:"name" msgpack:"name"`
	Value string `json:"value" msgpack:"value"`
}

type SetVar struct {
	Annotations   []interface{}       `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}         `json:"doc" msgpack:"doc"`
	Parameters    []interface{}       `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}       `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}       `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string              `json:"signature" msgpack:"signature"`
}

type SetVarParametersDoc struct {
	Name  string `json:"name" msgpack:"name"`
	Value string `json:"value" msgpack:"value"`
}

type SetVvar struct {
	Annotations   []interface{}        `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}          `json:"doc" msgpack:"doc"`
	Parameters    []interface{}        `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SetVvarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}        `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}        `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string               `json:"signature" msgpack:"signature"`
}

type SetVvarParametersDoc struct {
	Name  string `json:"name" msgpack:"name"`
	Value string `json:"value" msgpack:"value"`
}

type Strwidth struct {
	Annotations   []interface{}         `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}           `json:"doc" msgpack:"doc"`
	Parameters    []interface{}         `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc StrwidthParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string              `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}         `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                `json:"signature" msgpack:"signature"`
}

type StrwidthParametersDoc struct {
	Text string `json:"text" msgpack:"text"`
}

type Subscribe struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc SubscribeParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type SubscribeParametersDoc struct {
	Event string `json:"event" msgpack:"event"`
}

type TabpageDelVar struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc TabpageDelVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}              `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type TabpageDelVarParametersDoc struct {
	Name    string `json:"name" msgpack:"name"`
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
}

type TabpageGetNumber struct {
	Annotations   []interface{}                 `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                   `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                 `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc TabpageGetNumberParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                      `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                 `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                        `json:"signature" msgpack:"signature"`
}

type TabpageGetNumberParametersDoc struct {
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
}

type TabpageGetVar struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc TabpageGetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                   `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type TabpageGetVarParametersDoc struct {
	Name    string `json:"name" msgpack:"name"`
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
}

type TabpageGetWin struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc TabpageGetWinParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                   `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type TabpageGetWinParametersDoc struct {
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
}

type TabpageIsValid struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc TabpageIsValidParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
}

type TabpageIsValidParametersDoc struct {
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
}

type TabpageListWins struct {
	Annotations   []interface{}                `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                  `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc TabpageListWinsParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                     `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                       `json:"signature" msgpack:"signature"`
}

type TabpageListWinsParametersDoc struct {
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
}

type TabpageSetVar struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc TabpageSetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}              `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type TabpageSetVarParametersDoc struct {
	Name    string `json:"name" msgpack:"name"`
	Tabpage string `json:"tabpage" msgpack:"tabpage"`
	Value   string `json:"value" msgpack:"value"`
}

type UIAttach struct {
	Annotations   []interface{}         `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}           `json:"doc" msgpack:"doc"`
	Parameters    []interface{}         `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc UIAttachParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}         `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}         `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                `json:"signature" msgpack:"signature"`
}

type UIAttachParametersDoc struct {
	Height  string `json:"height" msgpack:"height"`
	Options string `json:"options" msgpack:"options"`
	Width   string `json:"width" msgpack:"width"`
}

type Uidetach struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type UipumSetHeight struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc UipumSetHeightParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
}

type UipumSetHeightParametersDoc struct {
	Height string `json:"height" msgpack:"height"`
}

type UISetOption struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type UITryResize struct {
	Annotations   []interface{} `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}   `json:"doc" msgpack:"doc"`
	Parameters    []interface{} `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc interface{}   `json:"parameters_doc,omitempty" msgpack:"parameters_doc,omitempty"`
	Return        []interface{} `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{} `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string        `json:"signature" msgpack:"signature"`
}

type UITryResizeGrid struct {
	Annotations   []interface{}                `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                  `json:"doc" msgpack:"doc"`
	Parameters    []interface{}                `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc UITryResizeGridParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}                `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}                `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                       `json:"signature" msgpack:"signature"`
}

type UITryResizeGridParametersDoc struct {
	Grid   string `json:"grid" msgpack:"grid"`
	Height string `json:"height" msgpack:"height"`
	Width  string `json:"width" msgpack:"width"`
}

type Unsubscribe struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc UnsubscribeParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}            `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type UnsubscribeParametersDoc struct {
	Event string `json:"event" msgpack:"event"`
}

type WinClose struct {
	Annotations   []interface{}         `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}           `json:"doc" msgpack:"doc"`
	Parameters    []interface{}         `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinCloseParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}         `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}         `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                `json:"signature" msgpack:"signature"`
}

type WinCloseParametersDoc struct {
	Force  string `json:"force" msgpack:"force"`
	Window string `json:"window" msgpack:"window"`
}

type WinDelVar struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinDelVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type WinDelVarParametersDoc struct {
	Name   string `json:"name" msgpack:"name"`
	Window string `json:"window" msgpack:"window"`
}

type WinGetBuf struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetBufParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type WinGetBufParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinGetConfig struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetConfigParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinGetConfigParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinGetCursor struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetCursorParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinGetCursorParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinGetHeight struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetHeightParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinGetHeightParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinGetNumber struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetNumberParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinGetNumberParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinGetOption struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetOptionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                  `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinGetOptionParametersDoc struct {
	Name   string `json:"name" msgpack:"name"`
	Window string `json:"window" msgpack:"window"`
}

type WinGetPosition struct {
	Annotations   []interface{}               `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                 `json:"doc" msgpack:"doc"`
	Parameters    []interface{}               `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetPositionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                    `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}               `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                      `json:"signature" msgpack:"signature"`
}

type WinGetPositionParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinGetTabpage struct {
	Annotations   []interface{}              `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}                `json:"doc" msgpack:"doc"`
	Parameters    []interface{}              `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetTabpageParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                   `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}              `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                     `json:"signature" msgpack:"signature"`
}

type WinGetTabpageParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinGetVar struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string               `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type WinGetVarParametersDoc struct {
	Name   string `json:"name" msgpack:"name"`
	Window string `json:"window" msgpack:"window"`
}

type WinGetWidth struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinGetWidthParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                 `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type WinGetWidthParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinIsValid struct {
	Annotations   []interface{}           `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}             `json:"doc" msgpack:"doc"`
	Parameters    []interface{}           `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinIsValidParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []string                `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}           `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                  `json:"signature" msgpack:"signature"`
}

type WinIsValidParametersDoc struct {
	Window string `json:"window" msgpack:"window"`
}

type WinSetBuf struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinSetBufParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type WinSetBufParametersDoc struct {
	Buffer string `json:"buffer" msgpack:"buffer"`
	Window string `json:"window" msgpack:"window"`
}

type WinSetConfig struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinSetConfigParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []string                  `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinSetConfigParametersDoc struct {
	Config string `json:"config" msgpack:"config"`
	Window string `json:"window" msgpack:"window"`
}

type WinSetCursor struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinSetCursorParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinSetCursorParametersDoc struct {
	Pos    string `json:"pos" msgpack:"pos"`
	Window string `json:"window" msgpack:"window"`
}

type WinSetHeight struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinSetHeightParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinSetHeightParametersDoc struct {
	Height string `json:"height" msgpack:"height"`
	Window string `json:"window" msgpack:"window"`
}

type WinSetOption struct {
	Annotations   []interface{}             `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}               `json:"doc" msgpack:"doc"`
	Parameters    []interface{}             `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinSetOptionParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}             `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}             `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                    `json:"signature" msgpack:"signature"`
}

type WinSetOptionParametersDoc struct {
	Name   string `json:"name" msgpack:"name"`
	Value  string `json:"value" msgpack:"value"`
	Window string `json:"window" msgpack:"window"`
}

type WinSetVar struct {
	Annotations   []interface{}          `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}            `json:"doc" msgpack:"doc"`
	Parameters    []interface{}          `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinSetVarParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}          `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}          `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                 `json:"signature" msgpack:"signature"`
}

type WinSetVarParametersDoc struct {
	Name   string `json:"name" msgpack:"name"`
	Value  string `json:"value" msgpack:"value"`
	Window string `json:"window" msgpack:"window"`
}

type WinSetWidth struct {
	Annotations   []interface{}            `json:"annotations,omitempty" msgpack:"annotations,omitempty"`
	Doc           interface{}              `json:"doc" msgpack:"doc"`
	Parameters    []interface{}            `json:"parameters,omitempty" msgpack:"parameters,omitempty"`
	ParametersDoc WinSetWidthParametersDoc `json:"parameters_doc" msgpack:"parameters_doc"`
	Return        []interface{}            `json:"return,omitempty" msgpack:"return,omitempty"`
	Seealso       []interface{}            `json:"seealso,omitempty" msgpack:"seealso,omitempty"`
	Signature     string                   `json:"signature" msgpack:"signature"`
}

type WinSetWidthParametersDoc struct {
	Width  string `json:"width" msgpack:"width"`
	Window string `json:"window" msgpack:"window"`
}
