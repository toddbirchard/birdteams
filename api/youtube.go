package api

import (
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
)

type YoutubeVideo struct {
	Title       string
	Description string
	Url         string
	Thumbnail   string
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
					Title:       videoTitle,
					Description: videoDescription,
					Thumbnail:   videoThumbnail,
					Url:         videoUrl,
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
		log.Fatalf("Error fetching videos: %v", err.Error())
	}

	channelResponse := getYoutubeChannel(service, "contentDetails")
	videos := getVideosFromChannel(service, channelResponse)
	return videos
}
