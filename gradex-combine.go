/*
 * Combine marked scripts from multiple markers, into one script per student.
 */

package main

import (
	pdfextract "github.com/georgekinnear/gradex-extract/pdfextract"
	"flag"
	"time"
	"os"
	"fmt"
	//"regexp"
)

func main() {

	var inputDir string
	flag.StringVar(&inputDir, "inputdir", "./", "path of the folder containing the PDF files to be processed (will also check sub-folders with 'marker' in their name")
	
	var markerOrder string
	flag.StringVar(&markerOrder, "markerorder", "", "comma separated list of marker initials, specifying the order to show them in the PDF (defaults to blank)")

	flag.Parse()

	
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		// inputDir does not exist
		fmt.Println(err)
		os.Exit(1)
	}
	
	report_time := time.Now().Format("2006-01-02-15-04-05")

	// Look at all PDFs in inputDir (including subdirectories)
	fmt.Println("Looking at input directory: ",inputDir)
	
	// Read the raw form values, and save them as a csv
	csv_path := fmt.Sprintf("%s/01_raw_form_values-%s.csv", inputDir, report_time)
	form_values := pdfextract.ReadFormsInDirectory(inputDir, csv_path)
	
	// Check the scripts are all from the same course
	coursecode := make(map[string]bool)
	for _, entry := range form_values {
		coursecode[entry.CourseCode] = true
	}
	if len(coursecode) != 1 {
		fmt.Println("Error - found scripts from multiple courses:",coursecode)
		os.Exit(1)
	}

	// Now summarise the marks and perform validation checks
	//csv_path = fmt.Sprintf("%s/00_marks_summary-%s.csv", inputDir, report_time)
	//pdf.ValidateMarking(form_values, parts, csv_path)

	err := mergePdf([]string{"test/marker1/B078068-mark.pdf", "test/marker2/B078068-mark.pdf"}, "output.pdf")
	if err != nil {
		fmt.Println(err)
	}
	
	pdfextract.PrettyPrintStruct(pdfextract.ReadFormFromPDF("output.pdf", false))
	
	os.Exit(1)

}
/*
func assemblePDFsInDirectory(inputDir string) {
	form_vals := []FormValues{}
	script_copies := make(map[string][]string)  // script_copies["B123456"] = []
	
	filename_examno, err := regexp.Compile("(B[0-9]{6})-.*.pdf")
	
	var num_scripts int
	filepath.Walk(inputDir, func(path string, f os.FileInfo, _ error) error {
		//if f.IsDir() && strings.Contains(f.Name(), "Moderation") { // TODO - check that this does not prevent us checking moderated marks!
		//	return filepath.SkipDir
		//}
		fmt.Println(f.Name())
		if !f.IsDir() {
			if filepath.Ext(f.Name()) != ".pdf" {
				return nil
			}
			proper_filename := filename_examno.MatchString(f.Name())
			if proper_filename {
				extracted_examno := filename_examno.FindStringSubmatch(f.Name())[1]
				vals_on_this_form := ReadFormFromPDF(path, true)
				// check that extracted_examno matches the one on the script!
				if vals_on_this_form[0].ExamNumber != extracted_examno {
					fmt.Println(" - Exam number mismatch: file",path,"has value",vals_on_this_form[0].ExamNumber)
				}
				
				form_vals = append(form_vals, vals_on_this_form...)
				num_scripts++
			} else {
				fmt.Println(" - Malformed filename: ", f.Name())
			}
		}
		return nil
	})
	
	
}
*/
