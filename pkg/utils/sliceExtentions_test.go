package utils

import (
	"github.com/stretchr/testify/assert"
	"storj/ditto/pkg/models"
	"testing"

	minio "github.com/minio/minio/cmd"
)

func TestListBucketsWithDifference(t *testing.T) {
	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "test 1",

			testFunc: func() {
				var mainSlice []minio.BucketInfo
				var mirrorSlice []minio.BucketInfo

				mainSlice = append(mainSlice, minio.BucketInfo{Name: "1"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "2"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "3"})

				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "4"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "2"})

				expectedSlice := []models.DiffModel{
					{ Name: "1", Diff: models.IN_MAIN },
					{ Name: "2", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME },
					{ Name: "3", Diff: models.IN_MAIN },
					{ Name: "4", Diff: models.IN_MIRROR },
				}

				result := ListBucketsWithDifference(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
					assert.Equal(t, expectedSlice[i].Diff, result[i].Diff)
				}
			},
		},
		{
			testName: "test 2",

			testFunc: func() {
				var mainSlice []minio.BucketInfo
				var mirrorSlice []minio.BucketInfo

				mainSlice = append(mainSlice, minio.BucketInfo{Name: "1"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "2"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "4"})

				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "4"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "2"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "2"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "2"})

				expectedSlice := []models.DiffModel{
					{ Name: "1", Diff: models.IN_MAIN },
					{ Name: "2", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME },
					{ Name: "3", Diff: models.IN_MAIN },
					{ Name: "4", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME },
				}

				result := ListBucketsWithDifference(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
					assert.Equal(t, expectedSlice[i].Diff, result[i].Diff)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}
}

func TestListObjectsWithDifference(t *testing.T) {
	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "test 1",

			testFunc: func() {
				var mainSlice []minio.ObjectInfo
				var mirrorSlice []minio.ObjectInfo

				//common
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "1", IsDir: true, Size: 1000, ContentType: "type1"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "1", IsDir: true, Size: 1000, ContentType: "type1"})


				//same name, diff size and type
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "4", IsDir: false, Size: 4000, ContentType: "type4"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "4", IsDir: false, Size: 4020, ContentType: "type4"})

				//only here
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "9", IsDir: true, Size: 7000, ContentType: "type7"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "7", IsDir: true, Size: 7000, ContentType: "type7"})


				expectedSlice := []models.DiffModel{
					{ Name: "1", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME | models.SIZE | models.CONTENT_TYPE | models.IS_DIR },
					{ Name: "4", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME | models.CONTENT_TYPE | models.IS_DIR },
					{ Name: "7", Diff: models.IN_MAIN },
					{ Name: "9", Diff: models.IN_MIRROR },
				}

				result := ListObjectsWithDifference(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
					assert.Equal(t, expectedSlice[i].Diff, result[i].Diff)
				}
			},
		},
		{
			testName: "test 2",

			testFunc: func() {
				var mainSlice []minio.ObjectInfo
				var mirrorSlice []minio.ObjectInfo

				//common
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "1", IsDir: true, Size: 1000, ContentType: "type1"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "2", IsDir: true, Size: 2000, ContentType: "type2"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "3", IsDir: true, Size: 3000, ContentType: "type3"})

				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "1", IsDir: true, Size: 1000, ContentType: "type1"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "2", IsDir: true, Size: 2000, ContentType: "type2"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "3", IsDir: true, Size: 3000, ContentType: "type3"})

				//same name
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "4", IsDir: false, Size: 4000, ContentType: "type4"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "5", IsDir: false, Size: 5000, ContentType: "type5"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "6", IsDir: false, Size: 6000, ContentType: "type6"})

				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "4", IsDir: false, Size: 4020, ContentType: "type4"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "5", IsDir: false, Size: 6000, ContentType: "type55"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "6", IsDir: true, Size: 7000, ContentType: "type6666"})

				//only here
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "7", IsDir: true, Size: 7000, ContentType: "type7"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "8", IsDir: true, Size: 8000, ContentType: "type8"})

				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "9", IsDir: true, Size: 7000, ContentType: "type7"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "10", IsDir: true, Size: 8000, ContentType: "type8"})

				expectedSlice := []models.DiffModel{
					{ Name: "1", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME | models.SIZE | models.CONTENT_TYPE | models.IS_DIR },
					{ Name: "2", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME | models.SIZE | models.CONTENT_TYPE | models.IS_DIR },
					{ Name: "3", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME | models.SIZE | models.CONTENT_TYPE | models.IS_DIR },

					{ Name: "4", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME | models.CONTENT_TYPE | models.IS_DIR },
					{ Name: "5", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME | models.IS_DIR },
					{ Name: "6", Diff: models.IN_MAIN | models.IN_MIRROR | models.NAME },

					{ Name: "7", Diff:  models.IN_MAIN  },
					{ Name: "8", Diff:  models.IN_MAIN  },
					{ Name: "9", Diff:  models.IN_MIRROR },
					{ Name: "10", Diff: models.IN_MIRROR },
				}


				result := ListObjectsWithDifference(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
					assert.Equal(t, expectedSlice[i].Diff, result[i].Diff)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}
}

func TestCombineBucketsDistinct(t *testing.T) {
	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "test 1",

			testFunc: func() {
				var mainSlice []minio.BucketInfo
				var mirrorSlice []minio.BucketInfo

				mainSlice = append(mainSlice, minio.BucketInfo{Name: "1"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "2"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "3"})

				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "1"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "4"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "2"})

				expectedSlice := []minio.BucketInfo{
					{Name: "1"},
					{Name: "2"},
					{Name: "3"},
					{Name: "4"},
				}

				result := CombineBucketsDistinct(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
				}
			},
		},
		{
			testName: "test 2",

			testFunc: func() {

				var mainSlice []minio.BucketInfo
				var mirrorSlice []minio.BucketInfo

				mainSlice = append(mainSlice, minio.BucketInfo{Name: "1"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "2"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.BucketInfo{Name: "5"})

				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "5"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "1"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "4"})
				mirrorSlice = append(mirrorSlice, minio.BucketInfo{Name: "2"})

				expectedSlice := []minio.BucketInfo{
					{Name: "1"},
					{Name: "2"},
					{Name: "3"},
					{Name: "5"},
					{Name: "4"},
				}

				result := CombineBucketsDistinct(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}
}

func TestCombineObjectsDistinct(t *testing.T) {
	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "test 1",

			testFunc: func() {
				var mainSlice []minio.ObjectInfo
				var mirrorSlice []minio.ObjectInfo

				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "1"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "2"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "3"})

				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "1"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "4"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "2"})

				expectedSlice := []minio.BucketInfo{
					{Name: "1"},
					{Name: "2"},
					{Name: "3"},
					{Name: "4"},
				}

				result := CombineObjectsDistinct(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
				}
			},
		},
		{
			testName: "test 2",

			testFunc: func() {

				var mainSlice []minio.ObjectInfo
				var mirrorSlice []minio.ObjectInfo

				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "1"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "2"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "3"})
				mainSlice = append(mainSlice, minio.ObjectInfo{Name: "5"})

				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "5"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "1"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "4"})
				mirrorSlice = append(mirrorSlice, minio.ObjectInfo{Name: "2"})

				expectedSlice := []minio.BucketInfo{
					{Name: "1"},
					{Name: "2"},
					{Name: "3"},
					{Name: "5"},
					{Name: "4"},
				}

				result := CombineObjectsDistinct(mainSlice, mirrorSlice)

				assert.Equal(t, len(expectedSlice), len(result))

				for i := 0; i < len(expectedSlice); i++ {
					assert.Equal(t, expectedSlice[i].Name, result[i].Name)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}
}