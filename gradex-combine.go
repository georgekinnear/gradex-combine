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
	"regexp"
	"path/filepath"
	"strings"
	"github.com/gocarina/gocsv"
)

func main() {

	var inputDir string
	flag.StringVar(&inputDir, "inputdir", "./", "path of the folder containing the PDF files to be processed (will also check sub-folders with 'marker' in their name")
	
	var outputDir string
	flag.StringVar(&outputDir, "outputdir", "scripts_combined", "path of the folder to receive the combined PDFs")
	
	var markerOrder string
	flag.StringVar(&markerOrder, "markerorder", "", "comma separated list of marker initials, specifying the order to show them in the PDF (defaults to blank)")

	flag.Parse()

	
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		// inputDir does not exist
		fmt.Println(err)
		os.Exit(1)
	}
	err := ensureDir(outputDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	report_time := time.Now().Format("2006-01-02-15-04-05")
	fmt.Println(report_time)
	csv_path := fmt.Sprintf("%s/00_combine-marks-%s.csv", inputDir, report_time)

	// Look at all PDFs in inputDir (including subdirectories)
	fmt.Println("Looking at input directory: ",inputDir)
	
	/*
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
	*/

	// Assemble a list of scripts, grouped by student
	filename_examno, err := regexp.Compile("(B[0-9]{6})-.*.pdf")
	
	scripts := make(map[string][]string)
	filepath.Walk(inputDir, func(path string, f os.FileInfo, _ error) error {
		// Avoid the moderation folder for now
		if f.IsDir() && strings.Contains(f.Name(), "Moderation") {
			return filepath.SkipDir
		}
		if f.IsDir() && strings.Contains(f.Name(), "scripts_combined") {
			return filepath.SkipDir
		}
		if !f.IsDir() {
			if filepath.Ext(f.Name()) != ".pdf" {
				return nil
			}
			proper_filename := filename_examno.MatchString(f.Name())
			if proper_filename {
				extracted_examno := filename_examno.FindStringSubmatch(f.Name())[1]
				if scripts[extracted_examno] == nil {
					scripts[extracted_examno] = []string{path}
				} else {
					scripts[extracted_examno] = append(scripts[extracted_examno], path)
				}
			} else {
				fmt.Println(" - Malformed filename: ", path)
			}
		}
		return nil
	})
	pdfextract.PrettyPrintStruct(scripts)
	
	
	
	form_vals_combined := []pdfextract.FormValues{}
	
	for ExamNo, script_set := range scripts {
		
		fmt.Println(ExamNo)
		output_file := outputDir+"/"+ExamNo+"-combinedmarks.pdf"
		
		// Merge the PDFs of this script
		err = mergePdf(script_set, output_file)
		if err != nil {
			fmt.Println(err)
		}
		
		// Read the details from the merged PDF and store them
		vals_on_new_form := pdfextract.ReadFormFromPDF(output_file, false)
		form_vals_combined = append(form_vals_combined, vals_on_new_form...)

	}
	
	// Put out a CSV summarising all the form values in the combined PDFs
	file, err := os.OpenFile(csv_path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	gocsv.MarshalFile(form_vals_combined, file)
	
	os.Exit(1)

	
	
	err = mergePdf([]string{"test/marker1/B078068-mark.pdf", "test/marker2/B078068-mark.pdf"}, "output.pdf")
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

// pr-pal @ https://stackoverflow.com/questions/37932551/mkdir-if-not-exists-using-golang
func ensureDir(dirName string) error {

	err := os.Mkdir(dirName, 0700) //probably umasked with 22 not 02

	os.Chmod(dirName, 0700)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}

}
