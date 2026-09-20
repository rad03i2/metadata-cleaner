# Metadata Cleaner

A small, dependency-free Go CLI for removing common privacy-sensitive metadata from **JPEG and PNG** images while preserving image payload data. It is designed for local, auditable workflows: no uploads, no telemetry, and no network access.

## Why it exists
Photos can contain EXIF/XMP/IPTC fields, comments, timestamps, software information, and sometimes location-related metadata. Metadata Cleaner creates a separate sanitized copy before an image is shared.

## Features
- Detects JPEG/PNG from file signatures rather than trusting extensions.
- JPEG: removes APP1 (EXIF/XMP), APP13 (IPTC/Photoshop) and COM comment segments.
- PNG: removes `tEXt`, `zTXt`, `iTXt`, `eXIf`, and `tIME` chunks.
- Never edits the source file in place.
- Refuses to replace an existing output unless `--overwrite` is explicit.
- Writes through a temporary file and syncs it before the final rename.
- Human-readable and JSON output.
- Standard-library-only implementation; no runtime third-party dependencies.

> **Preview guidance:** this is a CLI, so screenshots are optional. A useful repository preview is a terminal capture showing the input/output summary without personal file paths.

## Requirements
- Go 1.22+ to build from source.
- Windows, macOS, or Linux.

## Installation
```bash
git clone https://github.com/rad03i2/metadata-cleaner.git
cd metadata-cleaner
go build -o metadata-cleaner ./cmd/metadata-cleaner
```
On Windows the generated executable may be named `metadata-cleaner.exe`.

## Usage
```bash
# Creates photo.clean.jpg
./metadata-cleaner photo.jpg

# Choose a destination
./metadata-cleaner -output sanitized/photo.jpg photo.jpg

# Machine-readable result
./metadata-cleaner -json photo.png

# Explicitly replace an existing destination
./metadata-cleaner -output clean.jpg -overwrite photo.jpg

./metadata-cleaner -version
```
Exit code `0` means success, `1` means cleaning failed, and `2` means invalid CLI usage.

## Go API
```go
result, err := cleaner.CleanFile("photo.jpg", "photo.clean.jpg", false)
```
`Result` reports the detected format, byte counts, paths, and number of removed metadata blocks.

## Project structure
```text
cmd/metadata-cleaner/   CLI entry point
cleaner/                JPEG/PNG parser and cleaning engine
cleaner/*_test.go       automated safety/behavior tests
.github/workflows/      cross-platform CI
```

## Testing
```bash
go vet ./...
go test -race ./...
go build ./cmd/metadata-cleaner
```
CI runs these checks on Ubuntu, Windows, and macOS.

## Security & privacy
Processing is fully local. The tool does not transmit files. The original is preserved. Treat the cleaned copy as a new file and verify it visually before important publication. See [SECURITY.md](SECURITY.md) for reporting guidance.

## Limitations
- Supports JPEG and PNG only; WebP, HEIC/HEIF, TIFF, RAW, PDF, and video are not supported.
- It removes the documented metadata containers, not every conceivable steganographic or application-specific payload.
- JPEG ICC profiles (APP2) are preserved because they affect color rendering; APP0/other structural segments are preserved.
- PNG pixel content is not decoded/re-encoded; the tool removes selected ancillary chunks.
- This is a privacy utility, not a forensic anonymity guarantee.

## Optional roadmap
Potential future additions include WebP support, directory batch mode, dry-run metadata inspection, and reproducible release binaries. These are not implemented today.

## Contributing
Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md), add tests for behavior changes, and keep privacy-safe defaults.

## License
MIT — see [LICENSE](LICENSE).

## Author
**Radwan Abdulhadi Ahmed**  
**رضوان عبدالهادي أحمد**  
GitHub: **@rad03i2**

---

# منظّف البيانات الوصفية — العربية

أداة سطر أوامر صغيرة مكتوبة بلغة Go ومن دون اعتماديات خارجية، وظيفتها إزالة أشهر البيانات الوصفية الحساسة للخصوصية من صور **JPEG وPNG** مع الحفاظ على بيانات الصورة نفسها. تعمل محليًا بالكامل ولا ترفع الصور ولا تستخدم التتبع أو الشبكة.

## لماذا هذا المشروع؟
قد تحتوي الصور على EXIF وXMP وIPTC وتعليقات وتواريخ ومعلومات عن البرامج المستخدمة، وقد تتضمن بعض هذه البيانات معلومات مرتبطة بالموقع. تنشئ الأداة نسخة منفصلة ومنظفة قبل مشاركة الصورة.

## المزايا
- اكتشاف JPEG وPNG من توقيع الملف الحقيقي وليس من الامتداد فقط.
- في JPEG: إزالة APP1 الخاص غالبًا بـEXIF/XMP وAPP13 الخاص بـIPTC/Photoshop ومقاطع التعليقات COM.
- في PNG: إزالة `tEXt` و`zTXt` و`iTXt` و`eXIf` و`tIME`.
- عدم تعديل الملف الأصلي داخل مكانه نهائيًا.
- منع استبدال ملف ناتج موجود إلا عند استخدام `--overwrite` صراحةً.
- الكتابة أولًا إلى ملف مؤقت ومزامنته قبل اعتماد الناتج.
- مخرجات عادية أو JSON للأتمتة.
- لا توجد مكتبات تشغيل خارجية؛ يعتمد المحرك على مكتبة Go القياسية.

## المتطلبات والتثبيت
يتطلب Go 1.22 أو أحدث ويعمل على Windows وmacOS وLinux.
```bash
git clone https://github.com/rad03i2/metadata-cleaner.git
cd metadata-cleaner
go build -o metadata-cleaner ./cmd/metadata-cleaner
```

## الاستخدام
```bash
./metadata-cleaner photo.jpg
./metadata-cleaner -output sanitized/photo.jpg photo.jpg
./metadata-cleaner -json photo.png
./metadata-cleaner -output clean.jpg -overwrite photo.jpg
```
الرمز `0` للنجاح، و`1` لفشل التنظيف، و`2` لاستخدام غير صحيح للأوامر.

## الاختبار
```bash
go vet ./...
go test -race ./...
go build ./cmd/metadata-cleaner
```
وتنفذ GitHub Actions هذه الفحوصات على Linux وWindows وmacOS.

## البنية
`cmd/metadata-cleaner` يحتوي واجهة CLI، و`cleaner` يحتوي محرك تحليل وتنظيف JPEG/PNG، والاختبارات موجودة بجانب المحرك، و`.github/workflows` يحتوي إعداد CI.

## الأمان والخصوصية
كل المعالجة محلية ولا يتم إرسال الصور لأي خدمة. يبقى الأصل كما هو. يوصى بمعاينة النسخة الناتجة بصريًا قبل نشر الصور المهمة. تفاصيل الإبلاغ الأمني موجودة في [SECURITY.md](SECURITY.md).

## القيود
الأداة تدعم JPEG وPNG فقط حاليًا، ولا تدعم WebP أو HEIC أو TIFF أو RAW أو PDF أو الفيديو. وهي تزيل الحاويات المذكورة للبيانات الوصفية لكنها ليست ضمانًا ضد كل طرق إخفاء المعلومات. يتم الاحتفاظ بملفات تعريف الألوان ICC لأنها تؤثر في عرض الألوان، كما أن بيانات بكسلات PNG لا يعاد ترميزها.

## تطوير اختياري مستقبلًا
يمكن إضافة WebP، ووضع معالجة مجلد كامل، وفحص البيانات قبل الإزالة، وبناء إصدارات تنفيذية جاهزة. هذه الميزات غير موجودة حاليًا ولا يدعي المشروع وجودها.

## المساهمة والترخيص
راجع [CONTRIBUTING.md](CONTRIBUTING.md) قبل المساهمة، وأضف اختبارات لأي تغيير سلوكي. المشروع مرخص وفق MIT؛ راجع [LICENSE](LICENSE).

## المؤلف
**Radwan Abdulhadi Ahmed**  
**رضوان عبدالهادي أحمد**  
GitHub: **@rad03i2**
