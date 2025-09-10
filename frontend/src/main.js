function generateEventID() {
    return 'evt-' + Math.random().toString(36).slice(2) + Date.now();
  }
  
  document.getElementById('eventForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const action = document.getElementById('action').value;
    const tracking_id = document.getElementById('tracking_id').value;
    const user_id = document.getElementById('user_id').value;
    const event_id = generateEventID();
  
    // 构造表单数据
    const formData = new FormData();
    formData.append('tracking_id', tracking_id);
    formData.append('user_id', user_id);
    formData.append('event_id', event_id);
  
    // 发起请求
    let respBox = document.getElementById('respBox');
    respBox.textContent = '正在提交...';
  
    try {
      const resp = await fetch(`/event/${action}`, {
        method: 'POST',
        body: formData,
      });
      const resJson = await resp.json();
      respBox.textContent = '响应: ' + JSON.stringify(resJson, null, 2);
    } catch (err) {
      respBox.textContent = '请求失败: ' + err;
    }
  });