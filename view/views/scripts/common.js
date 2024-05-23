// ヘッダー部の高さ分だけコンテンツを下げる
// ヘッダー部分を取得
const header = document.querySelector('header');
const headerHeight = header.clientHeight;  // 高さをもらう

// bodyの上marginに設定
document.querySelector('main').style.paddingTop = headerHeight + 10 + 'px';