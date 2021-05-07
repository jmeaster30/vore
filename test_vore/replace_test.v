--replaces all instances of 'test-error' and 'test-fail' with 'test-success'
replace all "test-" = @a "error" or "fail" with @a "success"

--replaces all '-' with ': '
replace all '-' with ': '

--clears out all words in list
replace all in ['this', 'word', 'is', 'gone']