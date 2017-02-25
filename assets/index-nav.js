// Currently unused
//
// TODO:
//   - this works but steals _all_ keyboard events to the point that tabbing
//     and typing don't function
//   - this doesn't work in IE w/ arrows because of using event.key
var curIdx = -1;

function getChildren() {
  return $("#index-listing>tbody>tr");
}

function markActive(i) {
  c = getChildren();
  $($('td', c[i])[0]).addClass('active');
}

function markInactive(i) {
  c = getChildren();
  $($('td', c[i])[0]).removeClass('active');
}

function isActive(c, i) {
  return $($('td', c[i])[0]).hasClass('active');
}

function findActive(children) {
  if (curIdx != -1) { return; }

  c = getChildren();
  if (c.length == 0) { return; }

  for (var i = 0; i < c.length; i++) {
    if (isActive(c, i)) {
      curIdx = i;
      return;
    }
  }
}

function updateActive(up) {
  delta = 1;
  if (up) { delta = -1; }
  findActive();
  if (curIdx == -1) { return }

  c = getChildren();
  markInactive(curIdx);

  curIdx = curIdx + delta;
  if (curIdx < 0) { curIdx = c.length - 1; }
  if (curIdx >= c.length) { curIdx = 0; }
  markActive(curIdx);
}

function select() {
  findActive();
  if (curIdx == -1) { return; }
  liitem = getChildren()[curIdx];
  dest = $('a', liitem)[0].href;
  window.location.href = dest;
}

function keyToCmd(k) {
  k = k.toLowerCase();

  if (k == 'j' || k == 'arrowdown') {
    return 'down';
  }
  if (k == 'k' || k == 'arrowup') {
    return 'up';
  }
  if (k == '>' || k == 'enter') {
    return 'select';
  }
  return '';
}

/*
$(document).keypress(function(event) {
    cmd = keyToCmd(event.key);
    if (cmd == 'down') {
      updateActive(false);
    } else if (cmd == 'up') {
      updateActive(true);
    } else if (cmd == 'select') {
      select();
    } else {
      return false;
    }
    return true;
});
*/
