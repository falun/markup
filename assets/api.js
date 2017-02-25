function search(term, cb, opts) {
  if (opts === undefined) {
    opts = {}
  }

  opts['term'] = term;
  $.getJSON('/i/search', opts, cb)
}
