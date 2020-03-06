package main

import (
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	type args struct {
		version1 Version
		version2 Version
	}
	tests := []struct {
		name string
		args args
		want Version
	}{
		{
			args:args{
				version1: Version{major:1,minor:2,patch:0,preRelease:1,packageName:"b"},
				version2: Version{major:1,minor:1,patch:0,preRelease:1,packageName:"b"},
			},
			want:Version{major:1,minor:2,patch:0,preRelease:1,packageName:"b"},
		},
		{
			args:args{
				version1: Version{major:1,minor:2,patch:0,preRelease:1,packageName:"b"},
				version2: Version{major:2,minor:1,patch:0,preRelease:1,packageName:"b"},
			},
			want:Version{major:2,minor:1,patch:0,preRelease:1,packageName:"b"},
		},
		{
			args:args{
				version1: Version{major:1,minor:2,patch:0,preRelease:0,packageName:""},
				version2: Version{major:1,minor:2,patch:0,preRelease:5,packageName:"b"},
			},
			want:Version{major:1,minor:2,patch:0,preRelease:0,packageName:""},
		},
		//TODO: rc için test yazilacak
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.args.version1, tt.args.version2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name  string
		args  args
		want  Version
		want1 bool
	}{
		{
			args:args{tag:"1.0.3"},
			want:Version{major:1,minor:0,patch:3,packageName:"",preRelease:0},
			want1:false,
		},
		{
			args:args{tag:"1.0.3rc2"},
			want:Version{major:1,minor:0,patch:3,packageName:"rc",preRelease:2},
			want1:true,
		},
		{
			args:args{tag:"asdf"},
			want:Version{},
			want1:false,
		},
		{
			args:args{tag:"1.0.5"},
			want:Version{major:1,minor:0,patch:5},
			want1:false,
		},
		{
			args:args{tag:"1.2.3"},
			want:Version{major:1,minor:2,patch:3},
			want1:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseVersion(tt.args.tag)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}