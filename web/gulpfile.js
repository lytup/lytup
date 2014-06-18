var path = require('path');
var gulp = require('gulp');
var less = require('gulp-less');

gulp.task('less', function() {
  gulp.src('public/css/main.less')
    .pipe(less())
    .pipe(gulp.dest('public/css'));
});
