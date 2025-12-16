import { describe, it, expect } from "vitest";
import * as fc from "fast-check";
import { jsPDF } from "jspdf";

// Minimal valid JPEG base64 for testing (1x1 pixel)
const MINIMAL_JPEG =
  "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAn/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBEQCEAwEPwAB//9k=";

/**
 * Property 14: Image to PDF Conversion
 * For any valid image data, converting to PDF should produce a valid PDF file
 * **Validates: Requirements 7.2**
 */
describe("Property 14: Image to PDF Conversion", () => {
  it("should convert any image data to valid PDF", () => {
    fc.assert(
      fc.property(
        // Generate random dimensions for PDF placement
        fc.integer({ min: 10, max: 100 }),
        fc.integer({ min: 10, max: 100 }),
        (width: number, height: number) => {
          // Create PDF
          const pdf = new jsPDF({
            orientation: "portrait",
            unit: "mm",
            format: "a4",
          });

          // Add image to PDF with random dimensions
          pdf.addImage(MINIMAL_JPEG, "JPEG", 10, 10, width, height);

          // Get PDF output
          const pdfOutput = pdf.output("arraybuffer");

          // Verify PDF is valid (starts with %PDF)
          const pdfBytes = new Uint8Array(pdfOutput);
          const pdfHeader = String.fromCharCode(...pdfBytes.slice(0, 4));
          expect(pdfHeader).toBe("%PDF");

          // Verify PDF has content
          expect(pdfOutput.byteLength).toBeGreaterThan(0);

          return true;
        }
      ),
      { numRuns: 100 }
    );
  });
});

/**
 * Property 15: Multi-page PDF Combination
 * For any list of images, combining them should produce a PDF with exactly that many pages
 * **Validates: Requirements 7.3**
 */
describe("Property 15: Multi-page PDF Combination", () => {
  it("should create PDF with correct number of pages for any number of images", () => {
    fc.assert(
      fc.property(
        // Generate 1-10 images
        fc.integer({ min: 1, max: 10 }),
        (numImages: number) => {
          // Create multi-page PDF
          const pdf = new jsPDF({
            orientation: "portrait",
            unit: "mm",
            format: "a4",
          });

          for (let i = 0; i < numImages; i++) {
            if (i > 0) pdf.addPage();
            pdf.addImage(MINIMAL_JPEG, "JPEG", 10, 10, 50, 50);
          }

          // Verify page count
          const pageCount = pdf.getNumberOfPages();
          expect(pageCount).toBe(numImages);

          // Verify PDF is valid
          const pdfOutput = pdf.output("arraybuffer");
          expect(pdfOutput.byteLength).toBeGreaterThan(0);

          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should preserve page count when adding multiple images", () => {
    fc.assert(
      fc.property(
        fc.array(fc.integer({ min: 10, max: 50 }), {
          minLength: 2,
          maxLength: 5,
        }),
        (sizes: number[]) => {
          const pdf = new jsPDF({
            orientation: "portrait",
            unit: "mm",
            format: "a4",
          });

          sizes.forEach((size: number, index: number) => {
            if (index > 0) pdf.addPage();
            pdf.addImage(MINIMAL_JPEG, "JPEG", 10, 10, size, size);
          });

          // Verify correct number of pages
          expect(pdf.getNumberOfPages()).toBe(sizes.length);

          return true;
        }
      ),
      { numRuns: 100 }
    );
  });
});
