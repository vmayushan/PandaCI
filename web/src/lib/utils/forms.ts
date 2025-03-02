function formDataToJson<T extends Record<string, any> = Record<string, unknown>>(
	formData: FormData
): T {
	const obj: Record<string, any> = {};

	formData.forEach((value, key) => {
		// If the key already exists (for array values), append to the existing array
		if (obj[key]) {
			if (!Array.isArray(obj[key])) {
				obj[key] = [obj[key]];
			}
			obj[key].push(value);
		} else {
			obj[key] = value;
		}
	});

	return obj as T;
}

export function handleForm<T extends Record<string, any> = Record<string, unknown>>(
	event: SubmitEvent & {
		currentTarget: EventTarget & HTMLFormElement;
	}
) {
	event.preventDefault();
	return {
		data: formDataToJson<T>(new FormData(event.currentTarget)),
		value: (event.submitter as HTMLButtonElement | undefined)?.value
	};
}
