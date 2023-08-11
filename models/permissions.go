package models

const (
	ACTION_CAN_LOGIN_ADMIN = "can_login_admin"
	ACTION_CAN_LOGIN_USER  = "can_login_user"

	// ###### Product ######

	ACTION_PRODUCT_ADMIN_INFO   = "action_product_admin_info"
	ACTION_PRODUCT_ADMIN_LIST   = "action_product_admin_list"
	ACTION_PRODUCT_ADMIN_CREATE = "action_product_admin_create"
	ACTION_PRODUCT_ADMIN_UPDATE = "action_product_admin_update"
	ACTION_PRODUCT_ADMIN_DELETE = "action_product_admin_delete"

	// ###### Product ######

	// ###### Brand ######

	ACTION_BRAND_ADMIN_INFO   = "action_brand_admin_info"
	ACTION_BRAND_ADMIN_CREATE = "action_brand_admin_create"
	ACTION_BRAND_ADMIN_LIST   = "action_brand_admin_list"
	ACTION_BRAND_ADMIN_UPDATE = "action_brand_admin_update"
	ACTION_BRAND_ADMIN_DELETE = "action_brand_admin_delete"

	// ###### Brand ######

	// ###### Category ######

	ACTION_CATEGORY_ADMIN_INFO   = "action_category_admin_info"
	ACTION_CATEGORY_ADMIN_CREATE = "action_category_admin_create"
	ACTION_CATEGORY_ADMIN_LIST   = "action_category_admin_list"
	ACTION_CATEGORY_ADMIN_UPDATE = "action_category_admin_update"
	ACTION_CATEGORY_ADMIN_DELETE = "action_category_admin_delete"

	// ###### Category ######

	// ###### Role ######

	ACTION_ROLE_ADMIN_INFO        = "action_role_admin_info"
	ACTION_ROLE_ADMIN_CREATE      = "action_role_admin_create"
	ACTION_ROLE_ADMIN_LIST        = "action_role_admin_list"
	ACTION_ROLE_ADMIN_UPDATE      = "action_role_admin_update"
	ACTION_ROLE_ADMIN_DELETE      = "action_role_admin_delete"
	ACTION_ROLE_ADMIN_PERMISSIONS = "action_role_admin_permissions"

	// ###### Role ######

	// ###### Color ######

	ACTION_COLOR_ADMIN_INFO   = "action_color_admin_info"
	ACTION_COLOR_ADMIN_CREATE = "action_color_admin_create"
	ACTION_COLOR_ADMIN_LIST   = "action_color_admin_list"
	ACTION_COLOR_ADMIN_UPDATE = "action_color_admin_update"
	ACTION_COLOR_ADMIN_DELETE = "action_color_admin_delete"

	// ###### Color ######

	// ###### AppPic ######

	ACTION_APP_PIC_ADMIN_INFO   = "action_app_pic_admin_info"
	ACTION_APP_PIC_ADMIN_CREATE = "action_app_pic_admin_create"
	ACTION_APP_PIC_ADMIN_LIST   = "action_app_pic_admin_list"
	ACTION_APP_PIC_ADMIN_UPDATE = "action_app_pic_admin_update"
	ACTION_APP_PIC_ADMIN_DELETE = "action_app_pic_admin_delete"

	// ###### AppPic ######

	// ###### ProductFeatureCategory ######

	ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_INFO   = "action_product_feature_category_admin_info"
	ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_CREATE = "action_product_feature_category_admin_create"
	ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_LIST   = "action_product_feature_category_admin_list"
	ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_UPDATE = "action_product_feature_category_admin_update"
	ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_DELETE = "action_product_feature_category_admin_delete"

	// ###### ProductFeatureCategory ######

	// ###### ProductFeatureKey ######

	ACTION_PRODUCT_FEATURE_KEY_ADMIN_INFO   = "action_product_feature_key_admin_info"
	ACTION_PRODUCT_FEATURE_KEY_ADMIN_CREATE = "action_product_feature_key_admin_create"
	ACTION_PRODUCT_FEATURE_KEY_ADMIN_LIST   = "action_product_feature_key_admin_list"
	ACTION_PRODUCT_FEATURE_KEY_ADMIN_UPDATE = "action_product_feature_key_admin_update"
	ACTION_PRODUCT_FEATURE_KEY_ADMIN_DELETE = "action_product_feature_key_admin_delete"

	// ###### ProductFeatureKey ######

	// ###### ProductFeatureValue ######

	ACTION_PRODUCT_FEATURE_VALUE_ADMIN_INFO   = "action_product_feature_value_admin_info"
	ACTION_PRODUCT_FEATURE_VALUE_ADMIN_CREATE = "action_product_feature_value_admin_create"
	ACTION_PRODUCT_FEATURE_VALUE_ADMIN_LIST   = "action_product_feature_value_admin_list"
	ACTION_PRODUCT_FEATURE_VALUE_ADMIN_DELETE = "action_product_feature_value_admin_delete"

	// ###### ProductFeatureValue ######

	// ###### File ######

	ACTION_FILE_INFO            = "action_file_info"
	ACTION_FILE_CREATE          = "action_file_create"
	ACTION_FILE_LIST            = "action_file_list"
	ACTION_FILE_DELETE          = "action_file_delete"
	ACTION_FILE_CHANGE_PRIORITY = "action_file_change_priority"

	// ###### File - Systematic ######

	ACTION_FILE_SYSTEMATIC_INFO            = "action_file_systematic_info"
	ACTION_FILE_SYSTEMATIC_CREATE          = "action_file_systematic_create"
	ACTION_FILE_SYSTEMATIC_LIST            = "action_file_systematic_list"
	ACTION_FILE_SYSTEMATIC_DELETE          = "action_file_systematic_delete"
	ACTION_FILE_SYSTEMATIC_CHANGE_PRIORITY = "action_file_systematic_change_priority"

	// ###### File - systematic ######

	// ###### File - Product File Map ######

	ACTION_FILE_PRODUCT_INFO            = "action_file_product_info"
	ACTION_FILE_PRODUCT_CREATE          = "action_file_product_create"
	ACTION_FILE_PRODUCT_LIST            = "action_file_product_list"
	ACTION_FILE_PRODUCT_DELETE          = "action_file_product_delete"
	ACTION_FILE_PRODUCT_CHANGE_PRIORITY = "action_file_product_change_priority"

	// ###### File - Product File Map ######

	// ###### File - Brand ######

	ACTION_FILE_BRAND_INFO            = "action_file_brand_info"
	ACTION_FILE_BRAND_CREATE          = "action_file_brand_create"
	ACTION_FILE_BRAND_LIST            = "action_file_brand_list"
	ACTION_FILE_BRAND_DELETE          = "action_file_brand_delete"
	ACTION_FILE_BRAND_CHANGE_PRIORITY = "action_file_brand_change_priority"

	// ###### File - Brand ######

	// ###### File - AppPic ######

	ACTION_FILE_APP_PIC_INFO            = "action_file_app_pic_info"
	ACTION_FILE_APP_PIC_CREATE          = "action_file_app_pic_create"
	ACTION_FILE_APP_PIC_LIST            = "action_file_app_pic_list"
	ACTION_FILE_APP_PIC_DELETE          = "action_file_app_pic_delete"
	ACTION_FILE_APP_PIC_CHANGE_PRIORITY = "action_file_app_pic_change_priority"

	// ###### File - AppPic ######

	// ###### File ######

	// ###### VerificationCode ######

	ACTION_VERIFICATION_CODE_ADMIN_LIST = "action_verification_code_admin_list"

	// ###### VerificationCode ######

	// ###### Address ######

	ACTION_ADDRESS_ADMIN_LIST = "action_address_admin_list"

	// ###### Address ######
)

type Action struct {
	Name     string   `json:"name,omitempty"`
	Code     string   `json:"code,omitempty"`
	Children []Action `json:"children,omitempty"`
}

func GetPermissionsTree() []Action {
	actions := []Action{
		{
			Name: "Login Permissions",
			Children: []Action{
				{
					Name: "Ability to log in to admin panel with this role",
					Code: ACTION_CAN_LOGIN_ADMIN,
				},
				{
					Name: "Ability to log in to user panel with this role",
					Code: ACTION_CAN_LOGIN_USER,
				},
			},
		},
		{
			Name: "Categories",
			Code: "",
			Children: []Action{
				{
					Name: "Categories",
					Code: ACTION_CATEGORY_ADMIN_LIST,
				},
				{
					Name: "Create category",
					Code: ACTION_CATEGORY_ADMIN_CREATE,
				},
				{
					Name: "Update category",
					Code: ACTION_CATEGORY_ADMIN_UPDATE,
				},
				{
					Name: "Category information",
					Code: ACTION_CATEGORY_ADMIN_INFO,
				},
				{
					Name: "Delete category",
					Code: ACTION_CATEGORY_ADMIN_DELETE,
				},
			},
		},
		{
			Name: "Brands",
			Code: "",
			Children: []Action{
				{
					Name: "Brands",
					Code: ACTION_BRAND_ADMIN_LIST,
				},
				{
					Name: "Create brand",
					Code: ACTION_BRAND_ADMIN_CREATE,
				},
				{
					Name: "Update brand",
					Code: ACTION_BRAND_ADMIN_UPDATE,
				},
				{
					Name: "Brand information",
					Code: ACTION_BRAND_ADMIN_INFO,
				},
				{
					Name: "Delete brand",
					Code: ACTION_BRAND_ADMIN_DELETE,
				},
			},
		},
		{
			Name: "Roles",
			Code: "",
			Children: []Action{
				{
					Name: "Roles",
					Code: ACTION_ROLE_ADMIN_LIST,
				},
				{
					Name: "Create role",
					Code: ACTION_ROLE_ADMIN_CREATE,
				},
				{
					Name: "Update role",
					Code: ACTION_ROLE_ADMIN_UPDATE,
				},
				{
					Name: "Role information",
					Code: ACTION_ROLE_ADMIN_INFO,
				},
				{
					Name: "Delete role",
					Code: ACTION_ROLE_ADMIN_DELETE,
				},
				{
					Name: "Permissions",
					Code: ACTION_ROLE_ADMIN_PERMISSIONS,
				},
			},
		},
		{
			Name: "Colors",
			Code: "",
			Children: []Action{
				{
					Name: "Colors",
					Code: ACTION_COLOR_ADMIN_LIST,
				},
				{
					Name: "Create color",
					Code: ACTION_COLOR_ADMIN_CREATE,
				},
				{
					Name: "Update color",
					Code: ACTION_COLOR_ADMIN_UPDATE,
				},
				{
					Name: "Color information",
					Code: ACTION_COLOR_ADMIN_INFO,
				},
				{
					Name: "Delete color",
					Code: ACTION_COLOR_ADMIN_DELETE,
				},
			},
		},
		{
			Name: "App Pics",
			Code: "",
			Children: []Action{
				{
					Name: "App Pics",
					Code: ACTION_APP_PIC_ADMIN_LIST,
				},
				{
					Name: "Create app pic",
					Code: ACTION_APP_PIC_ADMIN_CREATE,
				},
				{
					Name: "Update app pic",
					Code: ACTION_APP_PIC_ADMIN_UPDATE,
				},
				{
					Name: "App pic information",
					Code: ACTION_APP_PIC_ADMIN_INFO,
				},
				{
					Name: "Delete app pic",
					Code: ACTION_APP_PIC_ADMIN_DELETE,
				},
			},
		},
		{
			Name: "Product Feature Categories",
			Code: "",
			Children: []Action{
				{
					Name: "Product Feature Categories",
					Code: ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_LIST,
				},
				{
					Name: "Create product feature category",
					Code: ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_CREATE,
				},
				{
					Name: "Update product feature category",
					Code: ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_UPDATE,
				},
				{
					Name: "Product feature category information",
					Code: ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_INFO,
				},
				{
					Name: "Delete product feature category",
					Code: ACTION_PRODUCT_FEATURE_CATEGORY_ADMIN_DELETE,
				},
			},
		},
		{
			Name: "Product Feature Keys",
			Code: "",
			Children: []Action{
				{
					Name: "Product Feature Keys",
					Code: ACTION_PRODUCT_FEATURE_KEY_ADMIN_LIST,
				},
				{
					Name: "Create product feature key",
					Code: ACTION_PRODUCT_FEATURE_KEY_ADMIN_CREATE,
				},
				{
					Name: "Update product feature key",
					Code: ACTION_PRODUCT_FEATURE_KEY_ADMIN_UPDATE,
				},
				{
					Name: "Product feature key information",
					Code: ACTION_PRODUCT_FEATURE_KEY_ADMIN_INFO,
				},
				{
					Name: "Delete product feature key",
					Code: ACTION_PRODUCT_FEATURE_KEY_ADMIN_DELETE,
				},
			},
		},
		{
			Name: "Product Feature Values",
			Code: "",
			Children: []Action{
				{
					Name: "Product Feature Values",
					Code: ACTION_PRODUCT_FEATURE_VALUE_ADMIN_LIST,
				},
				{
					Name: "Create product feature value",
					Code: ACTION_PRODUCT_FEATURE_VALUE_ADMIN_CREATE,
				},
				{
					Name: "Product feature value information",
					Code: ACTION_PRODUCT_FEATURE_VALUE_ADMIN_INFO,
				},
				{
					Name: "Delete product feature value",
					Code: ACTION_PRODUCT_FEATURE_VALUE_ADMIN_DELETE,
				},
			},
		},
		{
			Name: "Files",
			Code: "",
			Children: []Action{
				{
					Name: "Files",
					Code: ACTION_FILE_LIST,
				},
				{
					Name: "Create file",
					Code: ACTION_FILE_CREATE,
				},
				{
					Name: "File information",
					Code: ACTION_FILE_INFO,
				},
				{
					Name: "Delete file",
					Code: ACTION_FILE_DELETE,
				},
				{
					Name: "Change file priority",
					Code: ACTION_FILE_CHANGE_PRIORITY,
				},
				{
					Name: "Systematic Files",
					Code: "",
					Children: []Action{
						{
							Name: "Systematic Files",
							Code: ACTION_FILE_SYSTEMATIC_LIST,
						},
						{
							Name: "Create systematic file",
							Code: ACTION_FILE_SYSTEMATIC_CREATE,
						},
						{
							Name: "Systematic file information",
							Code: ACTION_FILE_SYSTEMATIC_INFO,
						},
						{
							Name: "Delete systematic file",
							Code: ACTION_FILE_SYSTEMATIC_DELETE,
						},
						{
							Name: "Change systematic file priority",
							Code: ACTION_FILE_SYSTEMATIC_CHANGE_PRIORITY,
						},
					},
				},
				{
					Name: "Product File Map",
					Code: "",
					Children: []Action{
						{
							Name: "Product File Map",
							Code: ACTION_FILE_PRODUCT_LIST,
						},
						{
							Name: "Create product file map",
							Code: ACTION_FILE_PRODUCT_CREATE,
						},
						{
							Name: "Product file map information",
							Code: ACTION_FILE_PRODUCT_INFO,
						},
						{
							Name: "Delete product file map",
							Code: ACTION_FILE_PRODUCT_DELETE,
						},
						{
							Name: "Change product file map priority",
							Code: ACTION_FILE_PRODUCT_CHANGE_PRIORITY,
						},
					},
				},
				{
					Name: "Brand Files",
					Code: "",
					Children: []Action{
						{
							Name: "Brand Files",
							Code: ACTION_FILE_BRAND_LIST,
						},
						{
							Name: "Create brand file",
							Code: ACTION_FILE_BRAND_CREATE,
						},
						{
							Name: "Brand file information",
							Code: ACTION_FILE_BRAND_INFO,
						},
						{
							Name: "Delete brand file",
							Code: ACTION_FILE_BRAND_DELETE,
						},
						{
							Name: "Change brand file priority",
							Code: ACTION_FILE_BRAND_CHANGE_PRIORITY,
						},
					},
				},
				{
					Name: "App Pic Files",
					Code: "",
					Children: []Action{
						{
							Name: "App Pic Files",
							Code: ACTION_FILE_APP_PIC_LIST,
						},
						{
							Name: "Create app pic file",
							Code: ACTION_FILE_APP_PIC_CREATE,
						},
						{
							Name: "App pic file information",
							Code: ACTION_FILE_APP_PIC_INFO,
						},
						{
							Name: "Delete app pic file",
							Code: ACTION_FILE_APP_PIC_DELETE,
						},
						{
							Name: "Change app pic file priority",
							Code: ACTION_FILE_APP_PIC_CHANGE_PRIORITY,
						},
					},
				},
			},
		},
		{
			Name: "Verification Codes",
			Code: "",
			Children: []Action{
				{
					Name: "Verification Codes",
					Code: ACTION_VERIFICATION_CODE_ADMIN_LIST,
				},
			},
		},
		{
			Name: "Address",
			Code: "",
			Children: []Action{
				{
					Name: "Address",
					Code: ACTION_ADDRESS_ADMIN_LIST,
				},
			},
		},
	}

	return actions
}
