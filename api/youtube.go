package api

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"strings"
)

// Go lib may be outdated:
// google.golang.org/api/googleapi/transport

// https://developers.google.com/youtube/v3/docs/search/list?apix_params=%7B%22part%22%3A%5B%22id%22%2C%22snippet%22%5D%2C%22channelId%22%3A%22UCZo-bCtK0xtclNm7HzksuWA%22%2C%22order%22%3A%22date%22%2C%22type%22%3A%5B%22video%22%5D%7D#go

type YoutubeVideo struct {
	Title          string
	Description    string
	Url            string
	Thumbnail      string
}

// Retrieve channel by ID
func getYoutubeChannel(service *youtube.Service, part string) *youtube.ChannelListResponse {
	channels := []string{part}
	call := service.Channels.List(channels)
	call = call.Id(os.Getenv("YOUTUBE_CHANNEL_ID"))
	response, err := call.Do()
	if err != nil {
		panic(err)
	}
	return response
}

// Retrieve items in the specified playlist
func playlistItemsList(service *youtube.Service, part string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	channels := []string{part}
	call := service.PlaylistItems.List(channels)
	call = call.PlaylistId(playlistId)
	call.MaxResults(100)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error when retrieving items from playlist: %v", err.Error())
	}
	return response
}

// Fetch videos from channel
func getVideosFromChannel(service *youtube.Service, response *youtube.ChannelListResponse) []YoutubeVideo {
	var videos []YoutubeVideo
	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads

		nextPageToken := ""
		for {
			// Retrieve next set of items in the playlist.
			playlistResponse := playlistItemsList(service, "snippet", playlistId, nextPageToken)

			for _, playlistItem := range playlistResponse.Items {
				videoTitle := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId
				videoDescription := strings.Split(playlistItem.Snippet.Description, "Bird teams")[0]
				videoThumbnail := playlistItem.Snippet.Thumbnails.High.Url
				videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)
				video := YoutubeVideo{
					Title: videoTitle,
					Description: videoDescription,
					Thumbnail: videoThumbnail,
					Url: videoUrl,
				}
				videos = append(videos, video)
				}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}
	return videos
}


func GetYoutubeVideos() []YoutubeVideo {
	ctx := context.Background()
	service, err := youtube.NewService(
		ctx,
		option.WithScopes(youtube.YoutubeReadonlyScope),
		option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))

	if err != nil {
		log.Fatalf( "Error fetching videos: %v", err.Error())
	}

	channelResponse := getYoutubeChannel(service, "contentDetails")
	videos := getVideosFromChannel(service, channelResponse)
	return videos
}