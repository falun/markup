var currentPage = 0;
var pageSize = 20;
var browseToken = '';

function ft(t) {
  return $('#ft-' + t)[0].checked;
}

function doSearch(newSearch) {
  if (newSearch === undefined || newSearch) {
    currentPage = 0;
  }

  var term = $("#search-term")[0].value;

  opts = {'ft':''};
  if (ft('md')) { opts['ft'] += 'md,'; }
  if (ft('go')) { opts['ft'] += 'go,'; }
  if (ft('js')) { opts['ft'] += 'js,'; }
  if ($('#ck-dir')[0].checked) { opts['include-dir'] = 'true'; }
  opts['page'] = currentPage;
  opts['page-size'] = pageSize;

  search(term, function(data) {
    populateResults(data, 'search-results');
  }, opts);
}

function populateResults(data, eleId) {
  var ele = $('#' + eleId)[0];
  ele.innerHTML = '';

  if (data['match'] === undefined) {
    console.log("Couldn't find match results");
    return;
  }

  results = data['match'];

  output = "<table class=\"table-striped table-condensed\">\n";
  output += "<tbody>\n";
  for (var idx in results) {
    path = results[idx];
    output += "<tr><td><a href='/" + browseToken + "/" + path + "'>" + path + "</a></td></tr>\n";
  }
  output += "</tbody>\n</table>\n";

  ele.innerHTML = output;

  cnt = data['count'];
  updatePager(currentPage, pageSize, cnt);
}

function updatePager(curPg, pgSz, cnt) {
  if (pgSz >= cnt) {
    hidePager();
    return;
  }

  var x = curPg + 1;
  var y = Math.ceil(cnt / pgSz);

  showPager();
  updateX(x);
  if (x == 1) { hidePrev(); } else { showPrev(); }

  updateY(y);
  if (x == y) { hideNext(); } else { showNext(); }
}

function showPrev() {
  $('#page-prev').show();
}

function hidePrev() {
  $('#page-prev').hide();
}

function showNext() {
  $('#page-next').show();
}

function hideNext() {
  $('#page-next').hide();
}

function showPager() {
  $('#page-div').show();
}

function hidePager() {
  $('#page-div').hide();
}

function updateX(x) {
  $('#x-of-y')[0].innerHTML = x;
}

function updateY(y) {
  $('#y-of-y')[0].innerHTML = y;
}

function prev() {
  currentPage--;
  doSearch(false);
}

function next() {
  currentPage++;
  doSearch(false);
}

function installSearch(dest, browse) {
  browseToken = browse;
  dest.innerHTML = '    <div class="modal fade" id="search-modal" role="dialog">\n' +
'  	  <div class="modal-dialog">\n' +
'\n' +
'  	    <!-- Modal content-->\n' +
'  	    <div class="modal-content">\n' +
'  	      <div class="modal-body">\n' +
'\n' +
'            <p>\n' +
'              <button\n' +
'                type="button"\n' +
'                style="float: right;"\n' +
'                class="btn btn-normal"\n' +
'                onclick="doSearch(true);"\n' +
'                tabindex="2">\n' +
'                Go\n' +
'              </button>\n' +
'              <b>Search</b> <input id="search-term" type="text" tabindex="1"></input><br/>\n' +
'              <b>File types</b>\n' +
'                .md: <input id="ft-md" type="checkbox" tabindex="3"></input>\n' +
'                .go: <input id="ft-go" type="checkbox" tabindex="4"></input>\n' +
'                .js: <input id="ft-js" type="checkbox" tabindex="5"></input>&nbsp;\n' +
'              <b>Include Path:</b> <input id="ck-dir" type="checkbox" tabindex="6"></input>\n' +
'              <br/>\n' +
'            </p>\n' +
'\n' +
'            <div id="page-div" class="row" style="width: 100%; display: none; margin-bottom: 0.5em;">\n' +
'              <div class="col-md-4 text-left"><a id="page-prev" href="#" onclick="prev();">Previous</a></div>\n' +
'              <div class="col-md-4 text-center">\n' +
'                Page <span id="x-of-y"></span> of <span id="y-of-y"></span>\n' +
'              </div>\n' +
'              <div style="float: right;" class="col-md-4 text-right"><a id="page-next" onclick="next();" href="#">Next</a></div>\n' +
'            </div>\n' +
'\n' +
'            <div id="search-results" class="row">\n' +
'            </div>\n' +
'\n' +
'  	      </div>\n' +
'\n' +
'  	      <div class="modal-footer">\n' +
'  	        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>\n' +
'  	      </div>\n' +
'  	    </div>\n' +
'\n' +
'  	  </div>\n' +
'  	</div>';

  $('#search-modal').on('shown.bs.modal', function() {
    $('#search-term').focus();
  });
}
