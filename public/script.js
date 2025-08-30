document.addEventListener('DOMContentLoaded', () => {

    const loadingMessage = document.getElementById('loading-message');
    const resultTable = document.getElementById('result-table');
    const resultBody = document.getElementById('result-body');

    async function fetchDomainStatus() {
        try {
            // ۱. ارسال درخواست به API بک‌اند ما
            const response = await fetch('/check');
            const results = await response.json();

            // ۲. نمایش جدول و پنهان کردن پیام "در حال بارگذاری"
            loadingMessage.classList.add('hidden');
            resultTable.classList.remove('hidden');
            resultBody.innerHTML = '';

            // ۳. حلقه زدن روی نتایج و ساختن ردیف‌های جدول
            results.forEach(result => {
                const row = document.createElement('tr');
                const statusText = result.isUp ? 'active' : 'not active';
                const statusClass = result.isUp ? 'status-up' : 'status-down';

                row.innerHTML = `
                    <td>${result.domain}</td>
                    <td class="${statusClass}">${statusText}</td>
                `;
                resultBody.appendChild(row);
            });

        } catch (error) {
            loadingMessage.innerText = 'خطا در برقراری ارتباط با سرور.';
            console.error('Error fetching data:', error);
        }
    }

    // اجرای تابع برای گرفتن داده‌ها
    fetchDomainStatus();
});