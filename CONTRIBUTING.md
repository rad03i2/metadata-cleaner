# Contributing / المساهمة

Thank you for improving Metadata Cleaner. Keep changes focused, dependency-light, and privacy-preserving.

1. Create a branch from `main`.
2. Add or update tests for behavioral changes.
3. Run `go vet ./...`, `go test -race ./...`, and `go build ./cmd/metadata-cleaner`.
4. Document user-visible behavior in both English and Arabic when applicable.
5. Open a focused pull request explaining the problem, solution, and validation.

Do not add telemetry, silent network calls, destructive in-place behavior, secrets, or test fixtures containing real personal metadata.

شكرًا للمساهمة. اجعل التغييرات محددة وتحافظ على الخصوصية، وأضف اختبارات لأي تغيير سلوكي. لا تضف تتبعًا أو اتصالات شبكة مخفية أو مفاتيح سرية أو صورًا حقيقية تحتوي معلومات شخصية ضمن الاختبارات.

Author / المؤلف: Radwan Abdulhadi Ahmed — رضوان عبدالهادي أحمد — @rad03i2
