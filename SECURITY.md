# Security Policy / سياسة الأمان

## Supported version
Security fixes target the latest code on `main`.

## Reporting
Please report suspected vulnerabilities privately through GitHub's security reporting features when available. Do not publish personal image metadata, exploit samples containing private information, credentials, or secrets in public issues.

Useful reports include the affected format, minimal synthetic reproduction steps, expected behavior, actual behavior, and platform/Go version.

## Security model
Metadata Cleaner processes local files only. It intentionally refuses in-place cleaning and refuses destination overwrite by default. It does not claim to provide forensic anonymity; callers remain responsible for verifying output appropriate to their threat model.

## العربية
تستهدف إصلاحات الأمان أحدث نسخة على `main`. عند اكتشاف مشكلة أمنية استخدم قنوات GitHub الأمنية الخاصة إن كانت متاحة، ولا تنشر بيانات وصفية شخصية أو أسرارًا ضمن Issue عامة. الأداة تعمل محليًا، ولا تعد ضمانًا للسرية الجنائية الكاملة؛ تحقق من الملف الناتج وفق مستوى الحماية المطلوب.

Maintainer / المسؤول: Radwan Abdulhadi Ahmed — رضوان عبدالهادي أحمد — @rad03i2
