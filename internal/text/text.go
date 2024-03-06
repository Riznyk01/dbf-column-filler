package text

const (
	StartMessage         = "Hello 👋 \nThis program is designed to modify column values in either a single file\nor a group of files that you drag and drop onto the executable file.\nSimply drag and drop one or more DBF files onto the executable,\nthen enter one or more pairs like col_name value, separated by a space.\nThe modified files will be saved in the 'changed' directory within\nthe same folder where the original files are located.\n\n"
	OutputFolder         = "changed"
	LogFileName          = "error_log.txt"
	FileExt              = "DBF"
	FailedToOpen         = "Failed to open log file"
	DropDBF              = "Please, drop only DBF files, try again"
	DropTheFiles         = "Please drop files on the executable file"
	EnterEven            = "You entered an odd number of parameters. Please, enter pairs.\n"
	DidntEnter           = "You didn't enter any parameters. Please, try again.\n"
	Working              = "Processing file:"
	CreatingFoldersError = "Error occurred while creating folders"
	FileSavedMessage     = "File %s has been successfully saved.\n"
	PanicMessage         = "Panic occurred:"
	SuccessMessage       = "Program completed successfully."
)
