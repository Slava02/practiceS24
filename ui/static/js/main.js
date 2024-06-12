var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') === window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

document.getElementById("createBtn").addEventListener("click", createForms);


function createForms() {
	const numForms = parseInt($('#num').val());
	if (isNaN(numForms) || numForms <= 0) {
		alert('Введите положительное число!');
		return;
	}

	$('.form-сontainer').empty();

	const title = `
						<div>
							<label>Название:</label>
							<input type="text" id="title" name='title'>
						</div><br>
						`;
	$('.form-сontainer').append(title);

	for (let i = 1; i <= numForms; i++) {
		const form = `
					<form class="myForm">
						<label>Точка №${i}</label>
						<br>
						<label for="x${i}">Координата x (float):</label>
						<input type="number" step="any" id="x${i}" name="x${i}" required>
						<br>
						<label for="y${i}">Координата y (float):</label>
						<input type="number" step="any" id="y${i}" name="y${i}" required>
						<br>
						<label for="z${i}">Координата z (float):</label>
						<input type="number" step="any" id="z${i}" name="z${i}" required>
						<br>
						<label for="mass${i}">Масса (float):</label>
						<input type="number" step="any" id="mass${i}" name="mass${i}" required>
						<br><br>
					</form>
                `;
		$('.form-сontainer').append(form);
	}
	const deleteIn = `
							<div>
								<label>Удалить через:</label>
								<input type='radio' name='expires' value='365' checked> Один год
								<input type='radio' name='expires' value='7'> Одну неделю
								<input type='radio' name='expires' value='1'> Один день
							</div>
							<br><br>
							`;
	$('.form-сontainer').append(deleteIn);

	const sendBtn = `<button type="submit" class="submitBtn">Создать вселенную</button>`
	$('.form-сontainer').append(sendBtn);

}

$(document).on('click', '.submitBtn', function(){
	const formData = [];
	const params = [];

	const title = $('#title').val();
	formData.push(title);

	$('.myForm').each(function(index, form) {
		const x = parseFloat($(form).find('input[name^="x"]').val());
		const y = parseFloat($(form).find('input[name^="y"]').val());
		const z = parseFloat($(form).find('input[name^="z"]').val());
		const mass = parseFloat($(form).find('input[name^="mass"]').val());

		if (isNaN(x) || isNaN(y) || isNaN(mass)) {
			alert('Пожалуйста, введите корректные значения для всех полей.');
			return;
		}

		params.push({coord: {x: x, y: y, z: z}, mass: mass});
	});

	formData.push(params);

	const expires = parseInt($('input[name="expires"]:checked').val());
	formData.push(expires);

	jQuery.ajax({
		url: '/universe/create',
		method: 'POST',
		data: JSON.stringify({title: formData[0], params: formData[1], expiresIn: formData[2]}),
		contentType: 'application/json',
		success: function(data, textStatus) {
			console.log("data Redirect: ", data.redirect)
			window.location = '/universe/view/' + data;
		},
		error: function(error) {
			console.error('Ошибка при отправке данных на сервер:', error);
		}
	});
});





